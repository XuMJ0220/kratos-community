package biz

import (
	"context"
	"strconv"

	"kratos-community/internal/conf"
	"kratos-community/internal/kafka"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 定义错误
var (
	ErrInternalServer   = errors.InternalServer("Err_INTERNAL_SERVER", "服务器出错")
	ErrArticleNotFound  = errors.NotFound("Article_Not_Found", "不存在该article")
	ErrAuthorizedUpdate = errors.Forbidden("FORBINDEN", "无权修改")
	ErrAuthorizedDelete = errors.Forbidden("FORBINDEN", "无权删除")
)

var (
	KAFKA_CREATE_ARTICLE_TOPIC = "article_created_topic"
)

// ContentRepo 与数据库交互的接口
type ContentRepo interface {
	CreateArtical(ctx context.Context, userid uint64, title, content string) (*Article, error)
	GetArticle(ctx context.Context, articleId uint64) (*Article, error)
	UpdateArticle(ctx context.Context, articleId uint64, title, content string) error
	DeleteArticle(ctx context.Context, articleId uint64) error
	DeleteArticleCache(ctx context.Context, articleId uint64) error
	// // CreateArticleInTx 在事务中创建文章且往outbox中插入一条需要往kafka生产的消息
	// CreateArticleInTx(ctx context.Context, tx *gorm.DB, article *Article) (*Article, error)
}

type ContentUsecase struct {
	repo        ContentRepo
	log         *log.Helper
	jwtSecret   string
	kafkaClient *kafka.KafkaClient
}

type Article struct {
	Id        uint64
	Title     string
	Content   string
	AuthorId  uint64
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
}

func NewContentUsecase(repo ContentRepo, logger log.Logger, jwtScret *conf.Auth, kafkaClient *kafka.KafkaClient) *ContentUsecase {
	return &ContentUsecase{repo: repo, log: log.NewHelper(logger), jwtSecret: jwtScret.JwtSecret, kafkaClient: kafkaClient}
}

func (uc *ContentUsecase) CreArticle(ctx context.Context, authorID uint64, title, content string) (*Article, error) {
	// 1.往数据库插入数据
	article, err := uc.repo.CreateArtical(ctx, authorID, title, content)
	if err != nil {
		return nil, err
	}

	// 往kafka中生产消息
	err = uc.kafkaClient.ProducerMessage(KAFKA_CREATE_ARTICLE_TOPIC, strconv.FormatUint(article.Id, 10), strconv.FormatUint(article.AuthorId, 10), "", "")
	if err != nil {
		uc.log.Errorf("ProducerMessage to kafka failed, err: %v", err)
	}

	// 2.返回结果
	return &Article{
		Id:        article.Id,
		Title:     article.Title,
		Content:   article.Content,
		AuthorId:  article.AuthorId,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}, nil
}

func (uc *ContentUsecase) GetArticle(ctx context.Context, articleId uint64) (*Article, error) {

	// 1.从数据库中查找
	article, err := uc.repo.GetArticle(ctx, articleId)
	if err != nil {
		return nil, err
	}

	// 2.返回结果
	return &Article{
		Id:        article.Id,
		Title:     article.Title,
		Content:   article.Content,
		AuthorId:  article.AuthorId,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}, nil
}

func (uc *ContentUsecase) UpdateArticle(ctx context.Context, articleId, authorId uint64, title, content string) (*Article, error) {
	// 进行授权检查
	article, err := uc.GetArticle(ctx, articleId)
	if err != nil {
		return nil, err
	}
	if article.AuthorId != authorId {
		return nil, ErrAuthorizedUpdate
	}

	// 更新数据库
	err = uc.repo.UpdateArticle(ctx, articleId, title, content)
	if err != nil {
		return nil, err
	}
	// 删除缓存
	err = uc.repo.DeleteArticleCache(ctx, articleId)
	// 获取最新数据
	article, _ = uc.GetArticle(ctx, articleId)
	return article, err
}

func (uc *ContentUsecase) DeleteArticle(ctx context.Context, articleId, authorId uint64) error {
	// 进行授权检查
	article, err := uc.GetArticle(ctx, articleId)
	if err != nil {
		return err
	}
	if article.AuthorId != authorId {
		return ErrAuthorizedDelete
	}
	// 先删除数据库
	if err := uc.repo.DeleteArticle(ctx, articleId); err != nil {
		// 如果是连数据库都删不了，就不用去删缓存了，等下次完整删除流程就行了
		return err
	}
	// 再删除缓存
	if err := uc.repo.DeleteArticleCache(ctx, articleId); err != nil {
		// 删除缓存的就打印一下日志就可以了
		uc.log.Errorf("DeleteArticleCache for id %d , error: %v", articleId, err)
	}
	return nil
}
