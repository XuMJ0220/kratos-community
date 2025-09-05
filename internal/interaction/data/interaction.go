package data

import (
	"context"
	"errors"
	"fmt"
	"kratos-community/internal/interaction/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-sql-driver/mysql"
)

type interactionRepo struct {
	data *Data
	log  *log.Helper
}

func NewInteractionRepo(data *Data, logger log.Logger) biz.InteractionRepo {
	return &interactionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateLike 实现创建点赞的接口
func (r *interactionRepo) CreateLike(ctx context.Context, like *biz.Like) error {
	// 将biz.Like 转换为data.Like(GORM模型)
	gormLike := &Like{
		UserID:    like.UserID,
		ArticleID: like.ArticleID,
	}
	// 插入数据库
	result := r.data.db1.WithContext(ctx).Create(gormLike)
	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		// 判断是不是唯一键冲突
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			r.log.Errorf("like_already_exist,delete_like from user_id : %d,article_id: %d failed", like.UserID, like.ArticleID)
			return biz.ErrLikeAlreadyExist
		}
		// 其他错误
		r.log.Errorf("CreateLike failed , error = %s", result.Error.Error())
		return biz.ErrInternalServer
	}
	return nil
}

// DeleteLike 实现删除点赞的接口
func (r *interactionRepo) DeleteLike(ctx context.Context, like *biz.Like) error {
	// 从数据库删除数据
	result := r.data.db1.WithContext(ctx).Where("user_id = ? AND article_id = ?", like.UserID, like.ArticleID).Delete(&Like{})
	// 出错了，但是不是因为找不到要删除的行
	if result.Error != nil {
		r.log.Errorf("DeleteLike failed , error = %s", result.Error.Error())
		return biz.ErrInternalServer
	}
	// 如果是找不到删除的行
	if result.RowsAffected == 0 {
		r.log.Errorf("like_not_found,delete_like from user_id : %d,article_id: %d failed", like.UserID, like.ArticleID)
		return biz.ErrLikeNotFound
	}
	return nil
}

// IncreLikeCount 实现点赞数递增的接口
func (r *interactionRepo) IncreLikeCount(ctx context.Context, articleId uint64) error {
	// 拼接Key
	key := fmt.Sprintf("article:like_count:%d", articleId)
	// 调用Redis的Incr方法
	return r.data.rdb1.Incr(key).Err()
}

// DecrLikeCount 实现点赞数递减的接口
func (r *interactionRepo) DecrLikeCount(ctx context.Context, articleId uint64) error {
	// 拼接Key
	key := fmt.Sprintf("article:like_count:%d", articleId)
	// 调用Redis的Decr方法
	return r.data.rdb1.Decr(key).Err()
}
