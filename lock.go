package lock

import (
	"context"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

type Lock struct {
	name    string
	owner   string
	seconds int
	redis   *redis.Client
}

type option func(lock *Lock)

func NewLock(redis *redis.Client, name string, opts ...option) *Lock {
	lock := &Lock{name: name, redis: redis}

	for _, opt := range opts {
		opt(lock)
	}

	if lock.owner == "" {
		lock.owner = randomString(16)
	}

	return lock
}

// WithOwner 锁的拥有者
func WithOwner(owner string) func(lock *Lock) {
	return func(lock *Lock) {
		lock.owner = owner
	}
}

// WithEX 锁的过期时间 单位 秒
func WithEX(seconds int) func(lock *Lock) {
	return func(lock *Lock) {
		lock.seconds = seconds
	}
}

// Acquire 加锁
func (lock *Lock) Acquire() bool {
	var val *redis.BoolCmd
	if lock.seconds > 0 {
		val = lock.redis.SetNX(context.Background(), lock.name, lock.owner, time.Duration(lock.seconds)*time.Second)
	} else {
		val = lock.redis.SetNX(context.Background(), lock.name, lock.owner, 0)
	}

	res, _ := val.Result()

	return res
}

// Release 释放自己加的锁
func (lock *Lock) Release() bool {
	eval := lock.redis.Eval(context.Background(), `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`, []string{lock.name}, lock.owner)
	b, err := eval.Bool()
	if err != redis.Nil {
		return false
	}
	return b
}

// ForceRelease 强制释放锁
func (lock *Lock) ForceRelease() {
	lock.redis.Del(context.Background(), lock.name)
}

// GetCurrentOwner 获取当前拥有者
func (lock *Lock) GetCurrentOwner() string {
	result, _ := lock.redis.Get(context.Background(), lock.name).Result()

	return result
}


const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // AllFilters 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}