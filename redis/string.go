package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type str struct {
	base
}

func NewString(key string) str {
	return str{base{key}}
}

// 合并了Set和SetEX
func (s str) Set(val any, exp time.Duration) error {
	if exp <= 0 {
		exp = redis.KeepTTL
	}
	return rdb.Set(ctx, s.key, val, exp).Err()
}

func (s str) SetNX(val any, exp time.Duration) (bool, error) {
	if exp <= 0 {
		exp = redis.KeepTTL
	}
	return rdb.SetNX(ctx, s.key, val, exp).Result()
}

// 注意的是，如果不存在则error为redis.Nil
func (s str) Get() (string, error) {
	v, err := rdb.Get(ctx, s.key).Result()
	return v, err
}

func (s str) IncrBy(incr int64) (int64, error) {
	return rdb.IncrBy(ctx, s.key, incr).Result()
}

func (s str) IncrByFloat(incr float64) (float64, error) {
	return rdb.IncrByFloat(ctx, s.key, incr).Result()
}
