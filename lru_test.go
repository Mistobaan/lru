package lru

import (
	"testing"
)

func TestGet(t *testing.T) {
	// TODO: 0 should be an error
	// TODO: Get from empty
	cache := NewLruCache()
	result := cache.Get("invalid")
	if result != nil {
		t.Error("nil expected")
	}
}
