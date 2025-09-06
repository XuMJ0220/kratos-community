package data

import (
	"context"
	"errors"
	"fmt"
	"kratos-community/internal/relation/biz"
	"strconv"

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
	key := fmt.Sprintf("following:%d", userId)
	return rr.data.rdb1.SAdd(key, followingId).Err()
}

// RemoveFollowing // 往userId的关注列表里移除followingId
func (rr *relationRepo) RemoveFollowing(ctx context.Context, userId, followingId uint64) error {
	key := fmt.Sprintf("following:%d", userId)
	return rr.data.rdb1.SRem(key, followingId).Err()
}

// AddFollower 往userId的粉丝列表里添加followerId(有人关注userId会触发)
func (rr *relationRepo) AddFollower(ctx context.Context, userId, followerId uint64) error {
	key := fmt.Sprintf("follower:%d", userId)
	return rr.data.rdb1.SAdd(key, followerId).Err()
}

// RemoveFollower 往userId的粉丝列表里移除followerId(有人取关userId会触发)
func (rr *relationRepo) RemoveFollower(ctx context.Context, userId, followerId uint64) error {
	key := fmt.Sprintf("follower:%d", userId)
	return rr.data.rdb1.SRem(key, followerId).Err()
}

// 列出userId的关注列表
// page 请求的第几页
// pageSize 每页的记录数
func (rr *relationRepo) ListFollowingIDs(ctx context.Context, userId, page, pageSize uint64) ([]uint64, int64, error) {
	key := fmt.Sprintf("following:%d", userId)

	// 获取关注的总数
	total, err := rr.data.rdb1.SCard(key).Result()
	if err != nil {
		rr.log.Errorf("ListFollowingIDs redis SCard error : %v", err)
		return nil, 0, biz.ErrInternalServer
	}
	if total == 0 {
		return []uint64{}, 0, nil
	}

	// 使用SMEMBERS命令，获取所有关注的ID(返回的是[]string)
	idStrs, err := rr.data.rdb1.SMembers(key).Result()
	if err != nil {
		rr.log.Errorf("ListFollowingIDs redis SMembers error : %v", err)
		return nil, 0, biz.ErrInternalServer
	}
	// 进行分页
	// 所有的范围是从0~(len(idStrs)-1)
	start := (page - 1) * pageSize
	end := start + pageSize // 这里不需要-1，因为后面用于切片操作，切片操作是左闭右开的
	// 处理边界情况
	if start >= uint64(len(idStrs)) {
		// 这种属于"没有记录"的情况，最佳就是返回一个200 OK和一个空数据列表
		// 因为如果用户点击了个超出限制的页数，返回个红色错误提示将体验非常差，应该提示“没有更多内容”这些就可以了
		return []uint64{}, total, nil // 请求的页码超过了总数据
	}
	if end > uint64(len(idStrs)) {
		end = uint64(len(idStrs))
	}
	// 切片处理，获取当前页的ID字符串
	paginatedIdStrs := idStrs[start:end]

	// 将字符串ID列表转换成uint64列表
	ids := make([]uint64, 0, len(paginatedIdStrs))
	for _, v := range paginatedIdStrs {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			rr.log.Errorf("ListFollowingIDs strconv.ParseUint error : %v", err)
			continue
		}
		ids = append(ids, id)
	}
	return ids, total, nil
}

func (rr *relationRepo) ListFollowerIDs(ctx context.Context, userId, page, pageSize uint64) ([]uint64, int64, error) {
	key := fmt.Sprintf("follower:%d", userId)

	// 获取关注的总数
	total, err := rr.data.rdb1.SCard(key).Result()
	if err != nil {
		rr.log.Errorf("ListFollowerIDs redis SCard error : %v", err)
		return nil, 0, biz.ErrInternalServer
	}
	if total == 0 {
		return []uint64{}, 0, nil
	}

	// 使用SMEMBERS命令，获取所有关注的ID(返回的是[]string)
	idStrs, err := rr.data.rdb1.SMembers(key).Result()
	if err != nil {
		rr.log.Errorf("ListFollowerIDs redis SMembers error : %v", err)
		return nil, 0, biz.ErrInternalServer
	}
	// 进行分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= uint64(len(idStrs)) {
		return []uint64{}, total, nil
	}
	if end > uint64(len(idStrs)) {
		end = uint64(len(idStrs))
	}
	// 切片处理，获取当前页的ID字符串
	paginatedIdStrs := idStrs[start:end]

	// 将字符串ID列表转换成uint64列表
	ids := make([]uint64, 0, len(paginatedIdStrs))
	for _, v := range paginatedIdStrs {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			rr.log.Errorf("ListFollowerIDs strconv.ParseUint error : %v", err)
			continue
		}
		ids = append(ids, id)
	}
	return ids, total, nil
}
