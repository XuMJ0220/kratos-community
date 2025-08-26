package biz

import (
	"context"

	"kratos-community/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 定义错误
var (
	ErrInternalServer = errors.InternalServer("Err_INTERNAL_SERVER", "服务器出错")
)

// ContentRepo 与数据库交互的接口
type ContentRepo interface {
	CreateArtical(ctx context.Context, userid uint64, title, content string) (*Article, error)
}

type ContentUsecase struct {
	repo      ContentRepo
	log       *log.Helper
	jwtSecret string
}

type Article struct {
	Id        uint64
	Title     string
	Content   string
	AuthorId  uint64
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
}

func NewContentUsecase(repo ContentRepo, logger log.Logger, jwtSecret *conf.Auth) *ContentUsecase {
	return &ContentUsecase{repo: repo, log: log.NewHelper(logger), jwtSecret: jwtSecret.JwtSecret}
}

func (uc *ContentUsecase) CreArticle(ctx context.Context, authorID uint64, title, content string) (*Article, error) {
	// 1.往数据库插入数据
	article,err:=uc.repo.CreateArtical(ctx,authorID,title,content)
	if err!=nil{
		return nil,err
	}
	// 2.返回结果
	return article, nil
}
