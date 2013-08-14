package lru

import (
	"testing"
)

func TestGet(t *testing.T) {
	// TODO: Get from empty
	cache := NewLruCache()
	result := cache.Get("invalid")
	if result != nil {
		t.Error("nil expected")
	}
}


func TestSet(t *testing.T) {
	cache := NewLruCache()
	exp := 1000
	result := cache.Set("key", exp)
	if result != nil {
		t.Error("nil expected")
	}
	value := cache.Get("key")
	if value != exp {
		t.Errorf("Expected %v got %v", exp, value)
	}
}
