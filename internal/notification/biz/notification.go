package biz

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	relationv1 "kratos-community/api/relation/v1"
)

var WG = sync.WaitGroup{}

// type NotificationRepo interface {
// 	// GetFollowers 方法用于获取一个用户的粉丝列表
// 	GetFollowers(ctx context.Context, userId uint64) (*relationv1.ListFollowersReply, error)
// }

type NotificationUseCase struct {
	//repo NotificationRepo
	log            *log.Helper
	relationClient relationv1.RelationClient
}

func NewNotificationUseCase(logger log.Logger, relationClient relationv1.RelationClient) *NotificationUseCase {
	return &NotificationUseCase{
		log:            log.NewHelper(logger),
		relationClient: relationClient,
	}
}

func (uc *NotificationUseCase) NotifyFollowers(ctx context.Context, followingId, articleId uint64)error{
	// 获取粉丝列表
	// 调用获取粉丝列表的grpc服务
	followers, err := uc.relationClient.ListFollowers(ctx, &relationv1.ListFollowersRequest{
		Id:       followingId,
		Page:     1,
		PageSize: 1000,
	})
	if err!=nil{
		uc.log.Errorf("获取粉丝列表失败,错误原因:%v", err)
		return err
	}
	// 发送通知，一共有followers.Total个粉丝
	// 每一页1000个粉丝，需要Total/1000+1页
	// 我们创建 Total/1000+1个goroutinie去通知
	// 每个goroutine获取一页粉丝，然后去通知
	allPages := followers.Total/1000 + 1 // 总页数
	var i uint64
	for i = 1; i <= allPages; i++ {
		WG.Add(1)

		go func(currentPage uint64) { // 使用传递参数
			defer WG.Done()
			fans, err := uc.relationClient.ListFollowers(ctx, &relationv1.ListFollowersRequest{
				Id:       followingId,
				Page:     currentPage,
				PageSize: 1000,
			})
			if err != nil {
				uc.log.Errorf("查找第 %d 页粉丝失败,每一页1000个粉丝,错误原因:%v", i, err)
			} else {
				// 暂时模拟通知粉丝
				for _, fan := range fans.Users {
					fmt.Printf("通知粉丝:%d\n 关注的用户:%d 更新了文章:%d\n", fan.Id, followingId, articleId)
				}
			}
		}(i) // 传递参数
	}
	WG.Wait()
	return nil
}
