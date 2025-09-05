package data

import (
	"context"
	"errors"
	"fmt"
	"kratos-community/internal/relation/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-sql-driver/mysql"
)

type relationRepo struct {
	data *Data
	log  *log.Helper
}

func NewRelationRepo(data *Data, logger log.Logger) biz.RelationRepo {
	return &relationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateRelation 创建关注关系
func (rr *relationRepo) CreateRelation(ctx context.Context, r *biz.Relation) error {
	gormRelation := &Relation{
		FollowerID:  r.FollowId,
		FollowingID: r.FollowingId,
	}

	result := rr.data.db1.WithContext(ctx).Create(gormRelation)

	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			return biz.ErrAlreadyFollowing
		}
		rr.log.Errorf("CreateRelation error : %v", result.Error)
		return biz.ErrInternalServer
	}
	return nil
}

// DeleteRelation 删除关注关系
func (rr *relationRepo) DeleteRelation(ctx context.Context, r *biz.Relation) error {
	result := rr.data.db1.WithContext(ctx).Where("follower_id = ? AND following_id = ?", r.FollowId, r.FollowingId).Delete(&Relation{})

	if result.Error != nil {
		rr.log.Errorf("DeleteRelation error : %v", result.Error)
		return biz.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return biz.ErrNotFollowing
	}
	return nil
}

// AddFollowing 往userId的关注列表里添加followingId
func (rr *relationRepo) AddFollowing(ctx context.Context, userId, followingId uint64) error {
	key := fmt.Sprintf("follower:%d", userId)
	return rr.data.rdb1.SAdd(key,followingId).Err()
}

// RemoveFollowing // 往userId的关注列表里移除followingId
func (rr *relationRepo) RemoveFollowing(ctx context.Context, userId, followingId uint64) error {
	key:=fmt.Sprintf("follower:%d", userId)
	return rr.data.rdb1.SRem(key,followingId).Err()
}

// AddFollower 往userId的粉丝列表里添加followerId(有人关注userId会触发)
func (rr *relationRepo) AddFollower(ctx context.Context, userId, followerId uint64) error {
	key:=fmt.Sprintf("following:%d", userId)
	return rr.data.rdb1.SAdd(key,followerId).Err()
}

// RemoveFollower 往userId的粉丝列表里移除followerId(有人取关userId会触发)
func (rr *relationRepo) RemoveFollower(ctx context.Context, userId, followerId uint64) error {
	key:=fmt.Sprintf("following:%d", userId)
	return rr.data.rdb1.SRem(key,followerId).Err()
}
