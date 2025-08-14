package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct { C *redis.Client }

func NewRedis(addr string, db int) *Redis {
	c := redis.NewClient(&redis.Options{Addr: addr, DB: db})
	return &Redis{C: c}
}

func (r *Redis) Ping(ctx context.Context) error { return r.C.Ping(ctx).Err() }

func (r *Redis) Close() error { return r.C.Close() }
