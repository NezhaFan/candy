package redisx

import (
	"errors"
	"math/rand"
	"time"
)

var (
	// 并发加锁重试次数
	retryTimes = 100
	// 并发加锁重试间隔
	retryInterval = 50 * time.Millisecond
	// 锁超时自动释放时间 (防止特殊情况未解锁)
	maxLockTime = 10 * time.Second
	// 随机数
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
	// 错误：加锁失败
	ErrRetryTimeout = errors.New("retry timeout")
	// 错误：上下文取消
	ErrContextCancel = errors.New("context cancel")
)

type Lock struct {
	base
	id int32
}

func NewLock(key string) *Lock {
	lock := &Lock{
		base: newBase("lock:" + key),
		id:   random.Int31(),
	}

	return lock
}

func (lock *Lock) TryLock() (bool, error) {
	return rdb.SetNX(lock.ctx, lock.key, lock.id, maxLockTime).Result()
}

func (lock *Lock) Lock() error {
	for i := 0; i <= retryTimes; i++ {
		select {
		// 上下文超时取消
		case <-lock.ctx.Done():
			return ErrContextCancel
		// 加锁
		default:
			ok, err := lock.TryLock()
			// 加锁出错
			if err != nil {
				return err
			}
			// 加锁成功
			if ok {
				return nil
			}

			time.Sleep(retryInterval)
		}
	}

	// 加锁超时
	return ErrRetryTimeout
}

const (
	// 删除。必须先匹配id值，防止A超时后，B马上获取到锁，A的解锁把B的锁删了
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

func (lock *Lock) Unlock() {
	rdb.Eval(lock.ctx, delCommand, []string{lock.key}, lock.id)
}
