package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"kratos-community/internal/conf"
	"kratos-community/internal/content/biz"
	"kratos-community/internal/pkg/redislock"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type contentRepo struct {
	data      *Data
	log       *log.Helper
	cacheMode *conf.CacheMode
	g         singleflight.Group
}

func NewContentRepo(data *Data, logger log.Logger, cacheMode *conf.CacheMode) biz.ContentRepo {
	return &contentRepo{
		data:      data,
		log:       log.NewHelper(logger),
		cacheMode: cacheMode,
		g:         singleflight.Group{},
	}
}

func (r *contentRepo) CreateArtical(ctx context.Context, userid uint64, title, content string) (*biz.Article, error) {
	// 插入数据
	article := Article{
		AuthorID: userid,
		Title:    title,
		Content:  content,
	}

	err := gorm.G[Article](r.data.db1).Create(ctx, &article)
	if err != nil {
		log.Errorf("CreateArtical error : %v", err)
		return nil, biz.ErrInternalServer
	}

	return &biz.Article{
		Id:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		AuthorId:  article.AuthorID,
		CreatedAt: timestamppb.New(article.CreatedAt), //  time.Time -> *timestamppb.Timestamp
		UpdatedAt: timestamppb.New(article.UpdatedAt),
	}, nil
}

func (c *contentRepo) GetArticle(ctx context.Context, articleId uint64) (*biz.Article, error) {
	// 1. 从redis中查询
	key := fmt.Sprintf("article:%d", articleId) // 拼接 Key
	result, err := c.data.rdb1.Get(key).Result()
	if err == nil {
		c.log.Infof("GetArticle from redis for key : %s", key)
		// 缓存命中
		var article biz.Article
		// 反序列化JSON字符串到结构体
		if err := json.Unmarshal([]byte(result), &article); err == nil {
			return &article, nil
		}
	}

	// 2.如果 err不是redis.Nil,说明redis服务器出错
	if err != redis.Nil {
		c.log.Errorf("redis Get error: %v", err)
	}
	// 3. redis 不命中,去mysql中查询
	c.log.Infof("redis miss, get article from mysql")
	// 用分布式锁防止缓存击穿
	lockKey := fmt.Sprintf("lock:article:%d", articleId)
	// 外面嵌套一层singleflght
	v, err, _ := c.g.Do(lockKey, func() (interface{}, error) {
		lock := redislock.NewRedisLock(c.data.rdb1, lockKey, 10*time.Second)
		locked, lockErr := lock.Lock(50*time.Millisecond, 3)
		// 获取锁失败
		if lockErr != nil || !locked {
			return nil, errors.New("concurrent lock failed, retry later")
		}
		defer lock.Unlock()

		article, err := gorm.G[Article](c.data.db1).Where("id = ?", articleId).First(ctx)
		if err != nil {
			c.log.Errorf("GetArticle: %v", err)         // 输出错误日志
			if errors.Is(err, gorm.ErrRecordNotFound) { // 不存在该行数据
				if c.cacheMode.CachePenetration == "1" {
					// 设置空缓存对象
					nilArtiStr, err := json.Marshal(biz.Article{})
					// 如果序列化失败，直接返回
					if err != nil {
						c.log.Errorf("json.Marshal error: %v", err)
						return nil, biz.ErrInternalServer
					}
					c.data.rdb1.Set(key, nilArtiStr, 1*time.Minute)
				}

				return nil, biz.ErrArticleNotFound
			} else { // 其他错误
				return nil, err
			}
		}
		if article.DeletedAt.Valid {
			return nil, biz.ErrArticleNotFound
		}

		bizArticle := biz.Article{
			Id:        articleId,
			Title:     article.Title,
			Content:   article.Content,
			AuthorId:  article.AuthorID,
			CreatedAt: timestamppb.New(article.CreatedAt),
			UpdatedAt: timestamppb.New(article.UpdatedAt),
		}

		// 4. 将mysql中的数据序列化成JSON字符串
		jsonData, err := json.Marshal(bizArticle)
		if err != nil {
			c.log.Errorf("json.Marshal error: %v", err)
			// 如果反序列化失败了，我们也只能提前返回了
			return &bizArticle, nil
		}

		// 5. 将JSON数据写入Redis缓存，并设置过期时间
		// 设置一个5分钟过期时间
		if err := c.data.rdb1.Set(key, jsonData, 5*time.Minute).Err(); err != nil {
			c.log.Errorf("redis Set error: %v", err)
		}
		return &bizArticle, nil
	})

	if err != nil {
		// 抢锁失败
		if err.Error() == "concurrent lock failed, retry later" {
			time.Sleep(time.Millisecond * 100)
			return c.GetArticle(ctx, articleId)
		}
		return nil, err
	}

	return v.(*biz.Article), nil
}

func (c *contentRepo) UpdateArticle(ctx context.Context, articleId uint64, title, content string) error {
	// 从数据库执行更新操作
	n, err := gorm.G[Article](c.data.db1).Where("id = ?", articleId).Updates(ctx, Article{Title: title, Content: content})
	// 如果是发生了错误
	if err != nil {
		c.log.Errorf("UpdateArticle: %v", err)
		return biz.ErrInternalServer
	}
	// 如果是查找不到文章
	if n == 0 {
		return biz.ErrArticleNotFound
	}
	// 来到了这里才表示更新成功了
	return nil
}

func (c *contentRepo) DeleteArticle(ctx context.Context, articleId uint64) error {
	_, err := gorm.G[Article](c.data.db1).Where("id = ?", articleId).Delete(ctx)
	if err != nil {
		c.log.Errorf("DeleteArticle: %v", err)
		return biz.ErrInternalServer
	}
	return nil
}

func (c *contentRepo) DeleteArticleCache(ctx context.Context, articleId uint64) error {
	// 拼接Key
	key := fmt.Sprintf("article:%d", articleId)
	// 从Redis中删除
	c.log.Infof("delete article from cache , key : %v", key)
	err := c.data.rdb1.Del(key).Err()
	return err
}
