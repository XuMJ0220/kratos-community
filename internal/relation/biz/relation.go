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
	ErrAlreadyFollowing = errors.Conflict("ALREADY_FOLLOWING", "已经关注,不需要再次关注")
	ErrNotFollowing     = errors.NotFound("NOT_FOLLOWING", "未关注,不能取消关注")
	ErrUserNotFound     = errors.NotFound("USER_NOT_FOUND", "用户不存在") // 用于跨服务调用检查
)

type Relation struct {
	FollowId    uint64 // 关注者id
	FollowingId uint64 // 被关注者id
}

type RelationRepo interface {
	// MySQL操作
	CreateRelation(ctx context.Context, r *Relation) error // 关注
	DeleteRelation(ctx context.Context, r *Relation) error // 取关
	// Redis操作
	AddFollowing(ctx context.Context, userId, followingId uint64) error    // 往userId的关注列表里添加followingId
	RemoveFollowing(ctx context.Context, userId, followingId uint64) error // 往userId的关注列表里移除followingId
	AddFollower(ctx context.Context, userId, followerId uint64) error      // 往userId的粉丝列表里添加followerId(有人关注userId会触发)
	RemoveFollower(ctx context.Context, userId, followerId uint64) error   // 往userId的粉丝列表里移除followerId(有人取关userId会触发)
}

type RelationUsecase struct {
	repo      RelationRepo
	log       *log.Helper
	jwtSecret string
}

func NewRelationUsecase(repo RelationRepo, logger log.Logger, jwtSecret *conf.Auth) *RelationUsecase {
	return &RelationUsecase{
		repo:      repo,
		log:       log.NewHelper(logger),
		jwtSecret: jwtSecret.JwtSecret,
	}
}

func (uc *RelationUsecase) FollowUser(ctx context.Context, followId, followingId uint64) error {
	// TODO: 在这里可以增加一次对 user-service 的 gRPC 调用，检查 followingID 是否存在

	err := uc.repo.CreateRelation(ctx, &Relation{
		FollowId:    followId,
		FollowingId: followingId,
	})
	if err != nil {
		return err
	}

	//更新 Redis 中的两个 Set
	// 我们允许这里的操作失败，只记录日志，保证核心功能可用
	err = uc.repo.AddFollowing(ctx, followId, followingId)
	if err != nil {
		uc.log.Errorf("FollowUser: AddFollowing failed: error: %v followId :%d followingId :%d ", err, followId, followingId)
	}
	err = uc.repo.AddFollower(ctx, followingId, followId)
	if err != nil {
		uc.log.Errorf("FollowUser: AddFollower failed: error: %v followId :%d followingId :%d ", err, followId, followingId)
	}
	return nil
}

func (uc *RelationUsecase) UnfollowUser(ctx context.Context, followId, followingId uint64) error {
	// TODO: 在这里可以增加一次对 user-service 的 gRPC 调用，检查 followingID  still exists

	err := uc.repo.DeleteRelation(ctx, &Relation{
		FollowId:    followId,
		FollowingId: followingId,
	})
	if err != nil {
		return err
	}

	// 更新 Redis 中的两个 Set
	// 我们允许这里的操作失败，只记录日志，保证核心功能可用
	err = uc.repo.RemoveFollowing(ctx, followId, followingId)
	if err != nil {
		uc.log.Errorf("UnfollowUser: RemoveFollowing failed: error: %v followId :%d followingId :%d ", err, followId, followingId)
	}
	err = uc.repo.RemoveFollower(ctx, followingId, followId)
	if err != nil {
		uc.log.Errorf("UnfollowUser: RemoveFollower failed: error: %v followId :%d followingId :%d ", err, followId, followingId)
	}
	return nil
}
