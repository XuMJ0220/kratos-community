package biz

import (
	"context"
	"kratos-community/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// 自定义错误
var (
	ErrInternalServer   = errors.InternalServer("Err_INTERNAL_SERVER", "服务器出错")
	ErrLikeAlreadyExist = errors.Conflict("LIKE_ALREADY_EXISTS", "已经点赞过了")
	ErrLikeNotFound = errors.NotFound("LIKE_NOT_FOUND", "还未点赞,无法取消点赞")
)

// Like 是点赞对象的业务对象
type Like struct {
	UserID    uint64
	ArticleID uint64
}

type InteractionRepo interface {
	CreateLike(ctx context.Context, like *Like) error // 创建点赞
	DeleteLike(ctx context.Context, like *Like) error // 删除点赞
	IncreLikeCount(ctx context.Context, articleId uint64) error // 增加点赞数量
	DecrLikeCount(ctx context.Context,articleId uint64) error // 减少点赞数量
}

type InteractionUsecase struct {
	repo      InteractionRepo
	log       *log.Helper
	jwtSecret string
}

func NewInteractionUsecase(repo InteractionRepo, logger log.Logger, jwtSecret *conf.Auth) *InteractionUsecase {
	return &InteractionUsecase{
		repo:      repo,
		log:       log.NewHelper(logger),
		jwtSecret: jwtSecret.JwtSecret,
	}
}

// LikeArticle 点赞文章
func (uc *InteractionUsecase) LikeArticle(ctx context.Context, userId, articleId uint64) error {
	// 创建点赞
	err:=uc.repo.CreateLike(ctx,&Like{UserID: userId, ArticleID: articleId})
	if err!=nil{
		return err
	}
	// 创建点赞成功，需要增加点赞数量
	err = uc.repo.IncreLikeCount(ctx,articleId)
	if err!=nil{
		// 如果更新 Redis 失败，我们只记录错误日志，不影响主流程的成功返回
		// 因为点赞记录本身已经成功写入数据库
		uc.log.Errorf("LikeArticle: IncrLikeCount failed after creating like record: %v", err)
	}
	return nil
}

// 
func (uc *InteractionUsecase) UnLikeArticle(ctx context.Context, userId, articleId uint64) error{
	// 删除点赞
	err:=uc.repo.DeleteLike(ctx,&Like{UserID: userId, ArticleID: articleId})
	if err!=nil{
		return err
	}
	// 删除点赞成功，需要减少点赞数量
	err = uc.repo.DecrLikeCount(ctx,articleId)
	if err!=nil{
		// 如果更新 Redis 失败，我们只记录错误日志，不影响主流程的成功返回
		// 因为点赞记录本身已经成功写入数据库
		uc.log.Errorf("UnLikeArticle: DecrLikeCount failed after deleting like record: %v", err)
	}
	return nil
}