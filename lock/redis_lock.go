package lock

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisLock
type RedisLock struct {
	client     *redis.Client
	key        string
	value      string
	xidFn      func() string
	expiration time.Duration
	mu         sync.Mutex
	stopRenew  chan struct{}
}

type Option func(*RedisLock)

func WithExpiration(exp time.Duration) Option {
	return func(l *RedisLock) {
		l.expiration = exp
	}
}
func WithXIDFn(fn func() string) Option {
	return func(rl *RedisLock) {
		rl.xidFn = fn
	}
}
func WithValue(value string) Option {
	return func(rl *RedisLock) {
		rl.value = value
	}
}
func xid() string {
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	return fmt.Sprintf("%s-%d-%d", hostname, pid, time.Now().Nanosecond())
}

// NewRedisLock 创建新锁,默认过期时间 10 s, 会在（过期时长）/2 的时刻自动续租
func NewRedisLock(client *redis.Client, key string, opts ...Option) *RedisLock {
	rs := &RedisLock{
		client:     client,
		key:        key,
		xidFn:      xid,
		expiration: time.Second * 10,
		stopRenew:  make(chan struct{}),
	}
	for _, opt := range opts {
		opt(rs)
	}
	if rs.value == "" {
		rs.value = rs.xidFn() //  value: 解锁和续约的时候，保证解除的锁是我自己添加的
	}
	return rs
}

// TryLock 尝试获取锁
// ok == false 时，说明锁已经存在
func (lock *RedisLock) TryLock() (ok bool, err error) {
	lock.mu.Lock()
	defer lock.mu.Unlock()

	ok, err = lock.client.SetNX(ctx, lock.key, lock.value, lock.expiration).Result()
	if err != nil {
		return false, err
	}
	if ok {
		go lock.autoRenewal()
	}
	return ok, nil
}

// Unlock 释放锁
func (lock *RedisLock) UnLock() error {
	lock.mu.Lock()
	defer lock.mu.Unlock()

	// 停止自动续租
	close(lock.stopRenew)

	// 使用 Lua 脚本确保只有持有锁的客户端可以释放锁
	script := `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end
`
	_, err := lock.client.Eval(ctx, script, []string{lock.key}, lock.value).Result()
	return err
}

// autoRenewal 自动续租
func (lock *RedisLock) autoRenewal() {
	ticker := time.NewTicker(lock.expiration / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lock.renew()
		case <-lock.stopRenew:
			return
		}
	}
}

// renew 续租
func (lock *RedisLock) renew() {
	lock.mu.Lock()
	defer lock.mu.Unlock()

	// 使用 Lua 脚本确保只有持有锁的客户端可以续租锁
	script := `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("expire", KEYS[1], ARGV[2])
else
	return 0
end
`
	lock.client.Eval(ctx, script, []string{lock.key}, lock.value, int(lock.expiration.Seconds())).Result()
}
