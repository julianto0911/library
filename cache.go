package library

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

type CacheConfiguration struct {
	URL      string
	Password string
	Prefix   string
	Port     string
}

func newRedisClient(url, port, password string, dbIndex int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     url + ":" + port,
		Password: password,
		DB:       dbIndex,
	})
}
func NewCache(cfg CacheConfiguration, expiracy int) Cache {
	return Cache{
		rdb:    newRedisClient(cfg.URL, cfg.Port, cfg.Password, 0),
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

// use this when init for ServiceContext, for local test
func MockCache(t *testing.T) Cache {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("fail init mock cache: %s", err)
	}

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	_, err = miniredis.Run()
	if err != nil {
		t.Fatalf("fail run mock cache: %s", err)
	}

	return Cache{
		rdb:    client,
		prefix: "",
	}
}
func LocalCache() (Cache, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return Cache{}, fmt.Errorf("fail init local cache: %s", err)
	}

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	_, err = miniredis.Run()
	if err != nil {
		return Cache{}, fmt.Errorf("fail  run local cache: %s", err)
	}

	return Cache{
		rdb:    client,
		prefix: "",
	}, nil
}
func (c *Cache) GetKeys() []string {
	var cursor uint64
	var result []string
	for {
		var keys []string
		var err error
		keys, cursor, err = c.rdb.Scan(ctxB, cursor, "", 0).Result()
		if err != nil {
			panic(err)
		}
		result = append(result, keys...)
		if cursor == 0 { // no more keys
			break
		}
	}

	return result
}

func (c *Cache) ClearKeys() {
	keys := c.GetKeys()
	pipe := c.rdb.Pipeline()
	for _, key := range keys {
		pipe.Del(ctxB, key)
	}
	pipe.Exec(ctxB)
}

func (c *Cache) Ping() bool {
	_, err := c.rdb.Ping(ctxB).Result()
	return err == nil
}
