package library

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

//use this when init for ServiceContext, for local test
func NewMockCacher(t *testing.T) Cacher {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("fail init mock cacher: %s", err)
	}

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	_, err = miniredis.Run()
	if err != nil {
		t.Fatalf("fail run mock cacher: %s", err)
	}

	return Cacher{
		rdb:      client,
		expiracy: time.Duration(10) * time.Second,
		prefix:   "",
	}
}

func TestCacher(t *testing.T) {
	r := NewMockCacher(t)
	err := r.Set("key", "val")
	assert.NoError(t, err)
	val, err := r.Get("key")
	assert.NoError(t, err)
	t.Logf("value :%s", val)
}
