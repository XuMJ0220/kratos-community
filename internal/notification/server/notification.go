package server

import (
	"context"
	"fmt"
	mykafka "kratos-community/internal/kafka"
	"kratos-community/internal/notification/biz"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/go-kratos/kratos/v2/log"
)

type NotificationServer struct {
	log *log.Helper
	uc  *biz.NotificationUseCase
	kc  *mykafka.KafkaConsumerClient
}

// NewNotificationServer 创建一个 NotificationServer
func NewNotificationServer(logger log.Logger, uc *biz.NotificationUseCase, kc *mykafka.KafkaConsumerClient) *NotificationServer {
	return &NotificationServer{
		log: log.NewHelper(logger),
		uc:  uc,
		kc:  kc,
	}
}

// Start 启动 NotificationServer
func (s *NotificationServer) Start(ctx context.Context) error {
	s.log.Info("start notification server...")

	// 启动一个 goroutine来执行消费逻辑，避免阻塞 Start 方法
	go func() {
		var run bool = true
		for run {
			select {
			// 监听上下文的取消信号，用于优雅退出
			case <-ctx.Done():
				run = false
				time.Sleep(15 * time.Second) // 延迟15秒退出
			default:
				// 从Kafka读取消息，超时时间为5秒
				msg, err := s.kc.Consumer.ReadMessage(1 * time.Second)
				if err != nil {
					// 正常的超时，说明暂时没有新消息，继续下一次循环即可
					if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
						fmt.Printf(".") // 打印一个点，表示正常的超时
						continue
					}
					// 真实的错误
					s.log.Errorf("!!! 消费错误: %v\n", err)
					run = false
					continue
				}
				// 调用biz层的业务来处理消息
				// 注意，这里我们新开一个goroutine去处理，避免耗时的业务阻塞下一条消息的消费
				go func() {
					followingId, _ := strconv.Atoi(string(msg.Key))
					articleId, _ := strconv.Atoi(string(msg.Value))
					if err := s.uc.NotifyFollowers(context.Background(), uint64(followingId), uint64(articleId)); err != nil {
						s.log.Errorf("!!! 处理消息错误: %v\n", err)
					}
				}()

			}
		}
	}()
	return nil
}

// Stop 停止 NotificationServer
func (s *NotificationServer) Stop(ctx context.Context) error {
	s.log.Info("stop notification server...")
	// 优雅关闭的逻辑已经在 NewKafkaConsumerClient 的 cleanup 函数中处理了
	return nil
}
