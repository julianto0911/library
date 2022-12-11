package library

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfiguration struct {
	URL      string
	Password string
	Prefix   string
}

func newRedisClient(url, password string, dbIndex int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       dbIndex,
	})
}
func NewCacher(cfg RedisConfiguration, expiracy int) Cacher {
	return Cacher{
		rdb:      newRedisClient(cfg.URL, cfg.Password, 0),
		expiracy: time.Duration(expiracy) * time.Second,
		prefix:   cfg.Prefix,
	}
}

func InjectMockCacher(rc *redis.Client, exp time.Duration, pref string) Cacher {
	return Cacher{
		rdb:      rc,
		expiracy: exp,
		prefix:   pref,
	}
}

type Cacher struct {
	rdb      *redis.Client
	expiracy time.Duration
	prefix   string
}

func (c *Cacher) Set(name string, value string) error {
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, c.expiracy).Err()
}

func (c *Cacher) SaveToken(name string, value string) error {
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, 24*60*60).Err()
}

func (c *Cacher) Get(name string) (string, error) {
	return c.rdb.Get(ctxB, c.prefix+"_"+name).Result()
}

func (c *Cacher) Delete(name string) error {
	return c.rdb.Del(ctxB, c.prefix+"_"+name).Err()
}
