package library

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

// use this when init for ServiceContext, for local test
func NewMockCache(t *testing.T) Cache {
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

func TestCache(t *testing.T) {
	r := NewMockCache(t)
	err := r.Set("key", "val", 1*time.Minute)
	assert.NoError(t, err)
	val, err := r.Get("key")
	assert.NoError(t, err)
	t.Logf("value :%s", val)
}
