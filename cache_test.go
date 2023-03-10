package library

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	r := MockCache(t)
	err := r.Set("key", "val", 1*time.Minute)
	assert.NoError(t, err)
	val, err := r.Get("key")
	assert.NoError(t, err)
	t.Logf("value :%s", val)
}
