package kafka

import (
	"kratos-community/internal/conf"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// KafkaClient 生产者客户端
type KafkaClient struct {
	producer *kafka.Producer
	log      *log.Helper
}

// KafkaConsumerClient 消费者客户端
type KafkaConsumerClient struct {
	Consumer *kafka.Consumer
	log      *log.Helper
}

var ProviderSet = wire.NewSet(NewKafkaClient, NewKafkaConsumerClient)

func NewKafkaClient(kaf *conf.Kafka, logger log.Logger) (*KafkaClient, func(), error) {
	// 1. 创建生产者配置
	config := &kafka.ConfigMap{
		"bootstrap.servers":  kaf.Bootstrap.Servers,
		"enable.idempotence": kaf.Enable.Idepotence,
		"acks":               kaf.Acks,
		"retries":            int(kaf.Retries),
	}

	log := log.NewHelper(logger)

	// 2. 创建生产者实例
	p, err := kafka.NewProducer(config)
	if err != nil {
		log.Errorf("kafka生存者实例创建失败, err: %v", err)
		return nil, nil, err
	}

	// 增加一个goroutine来异步处理交付报告
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("!!! 消息发送失败: %v\n", ev.TopicPartition.Error)
				} else {
					log.Infof(">>> 消息发送成功: Topic=%s, Partition=[%d], Offset=%v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	// 用于关闭生产者实例
	cleanup := func() {
		defer p.Close() // 关闭生产者实例

		// 等待消息发送完成,最多等待10秒
		// Flush 等待所有消息都被发送完毕，或者超时就返回
		unflushed_messages := p.Flush(15 * 1000) // 最多等待15秒
		if unflushed_messages > 0 {
			log.Errorf("有 %d 条消息未发送成功", unflushed_messages)
		} else {
			log.Info("消息发送成功")
		}
	}

	return &KafkaClient{
		producer: p,
		log:      log,
	}, cleanup, nil
}

func (k *KafkaClient) ProducerMessage(topic, key, value, headersKey, headersValue string) error {
	// 准备要生存的消息
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(value),
		Headers: []kafka.Header{
			{Key: headersKey, Value: []byte(headersValue)},
		},
	}

	// 往topic生产消息
	// 这是一个异步调用
	// Sender 线程会在后台批量发送消息
	err := k.producer.Produce(message, nil)
	if err != nil {
		k.log.Errorf("kafka producer send message failed, err: %v", err)
		return err
	}

	return nil
}

func NewKafkaConsumerClient(kaf *conf.Kafka, logger log.Logger) (*KafkaConsumerClient, func(), error) {
	// 1. 创建消费者配置
	config := &kafka.ConfigMap{
		"bootstrap.servers": kaf.Bootstrap.Servers,
		// 消费者组ID，所有使用相同group.id的消费者实例都属于同一组
		"group.id": kaf.Group.Id,
		// "earliest": 从最早的消息开始
		// "latest": 从最新的消息开始
		"auto.offset.reset": kaf.Auto.Offset.Reset_,
		// 是否自动提交位移，我们手动控制的时候是最可靠的
		"enable.auto.commit": kaf.Enable.Auto.Commit,
	}
	// 2. 创建消费者实例
	c, err := kafka.NewConsumer(config)
	if err != nil {
		log.Errorf("kafka消费者实例创建失败, err: %v", err)
		return nil, nil, err
	}

	// 3. 订阅主题
	err = c.SubscribeTopics(kaf.SubTopics,nil)
	if err!=nil{
		log.Errorf("kafka消费者订阅主题失败, err: %v", err)
		return nil, nil, err
	}
	
	cleanup := func() {
		log.Info("closing kafka consumer")
		c.Close()
	}

	return &KafkaConsumerClient{
		Consumer: c,
		log:      log.NewHelper(logger),
	}, cleanup, nil
}
