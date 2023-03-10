package library

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheConfiguration struct {
	URL      string
	Password string
	Prefix   string
	Port     string
}

func newRedisClient(url, password string, dbIndex int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       dbIndex,
	})
}
func NewCache(cfg CacheConfiguration, expiracy int) Cache {
	return Cache{
		rdb:    newRedisClient(cfg.URL, cfg.Password, 0),
		prefix: cfg.Prefix,
	}
}

type Cache struct {
	rdb    *redis.Client
	prefix string
}

func (c *Cache) Set(name string, value string, tm time.Duration) error {
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, tm).Err()
}

func (c *Cache) SaveToken(name string, value string, tm time.Duration) error {
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, tm).Err()
}

func (c *Cache) Get(name string) (string, error) {
	return c.rdb.Get(ctxB, c.prefix+"_"+name).Result()
}

func (c *Cache) Delete(name string) error {
	return c.rdb.Del(ctxB, c.prefix+"_"+name).Err()
}
