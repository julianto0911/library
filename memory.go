package library

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Memory struct {
	rdb    *redis.Client
	prefix string
}

func NewMemory(cfg RedisConfiguration) Cacher {
	return Cacher{
		rdb:    newRedisClient(cfg.URL, cfg.Password, 0),
		prefix: cfg.Prefix,
	}
}

func InjectMockMemory(rc *redis.Client, pref string) Memory {
	return Memory{
		rdb:    rc,
		prefix: pref,
	}
}

func (c *Memory) Set(name string, value string, expiracy time.Duration) error {
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, expiracy).Err()
}

func (c *Memory) Get(name string) (string, error) {
	return c.rdb.Get(ctxB, c.prefix+"_"+name).Result()
}

func (c *Memory) Delete(name string) error {
	return c.rdb.Del(ctxB, c.prefix+"_"+name).Err()
}
