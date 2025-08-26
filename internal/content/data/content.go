package data

import (
	"context"
	"kratos-community/internal/content/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type contentRepo struct {
	data *Data
	log  *log.Helper
}

func NewContentRepo(data *Data, logger log.Logger) biz.ContentRepo {
	return &contentRepo{
		data: data,
		log:  log.NewHelper(logger),
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
