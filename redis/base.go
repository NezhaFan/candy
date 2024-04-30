package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// 文档 https://redis.uptrace.dev/zh/guide/

var (
	rdb      *redis.Client
	ctx      = context.Background()
	longtime = time.Hour * 24 * 30
)

const (
	OK  = "OK"
	Nil = redis.Nil
)

func Connect(addr string, auth string) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: auth,
		DB:       0,
	})
	return rdb.Ping(ctx).Err()
}

type base struct {
	key string
}

func (b base) Expire(exp time.Duration) (bool, error) {
	return rdb.Expire(ctx, b.key, exp).Result()
}

func (b base) Exists() (bool, error) {
	i, err := rdb.Exists(ctx, b.key).Result()
	return i == 1, err
}

func (b base) Del() bool {
	return rdb.Del(ctx, b.key).Val() == 1
}

func (b base) TTL() time.Duration {
	return rdb.TTL(ctx, b.key).Val()
}

func (b base) Do(args ...any) (any, error) {
	args = append([]any{b.key}, args...)
	return rdb.Do(ctx, args...).Result()
}
func Pipeline() redis.Pipeliner {
	pipe := rdb.Pipeline()
	return pipe
}
