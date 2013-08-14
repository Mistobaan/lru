package lru

import (
	"testing"
)


// states of cache:
// empty, one, multiple, max
type Case struct {
	Name string
	InitialState []interface{}
	Result []struct {
		key string
		value interface{}
	}
	Message string
}


func TestGetEmpty(t *testing.T) {
	cache := NewLruCache()
	result := cache.Get("invalid")
	if result != nil {
		t.Error("nil expected")
	}
}

func TestSetEmpty(t *testing.T) {
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
