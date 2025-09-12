package redislock

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// RedisLock 是一个基于Redis的分布式锁
type RedisLock struct {
	client     *redis.Client
	key        string        // 锁的 key
	value      string        // 锁的 value 用于识别锁的拥有者
	expiration time.Duration // 锁的过期时间
}

// NewRedisLock 创建一个RedisLock实例
func NewRedisLock(client *redis.Client, key string, expiration time.Duration) *RedisLock {
	return &RedisLock{
		client:     client,
		key:        key,
		value:      uuid.NewString(), // 生成一个唯一的UUID作为锁的值
		expiration: expiration,
	}
}

// Lock 尝试获取锁（带重试机制）
// ctx: 上下文
// retryAfter: 每次重试的间隔
// maxRetries: 最大重试次数
func (l *RedisLock) Lock(retryAfter time.Duration, maxRetries int) (bool, error) {
	for i := 0; i < maxRetries; i++ {
		// 尝试使用SET NX PX 原子命令来获取锁
		ok, err := l.client.SetNX(l.key, l.value, l.expiration).Result()
		// 如果出错，返回错误
		if err != nil {
			return false, err
		}
		// 如果获取成功，返回true
		if ok {
			return true, nil
		}
		// 如果获取失败，等待一段时间后重试
		time.Sleep(retryAfter)
	}
	// 超过了最大重试次数，返回失败
	return false, nil
}

func (l *RedisLock) Unlock() error {
	// 这个Lua脚本会先GET key的值，如果和我们传入的value一致，再删除key
	script := `
	if redis.call("get",KEYS[1]) == ARGV[1] then
		return redis.call("del",KEYS[1])
	else
		return 0
	end
	`
	// 执行脚本
	_, err := l.client.Eval(script, []string{l.key}, l.value).Result()
	return err
}
