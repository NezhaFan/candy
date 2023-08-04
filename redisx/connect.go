package redisx

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

const (
	OK = "OK"
)

var (
	emptyCtx = context.Background()
)

var (
	rdb *redis.Client
)
