package lru

import (
	"testing"
)

// states of cache:
// empty, one, multiple, max
type Case struct {
	Name         string
	InitialState []interface{}
	Result       []struct {
		key   string
		value interface{}
	}
	Message string
}

func TestGetEmpty(t *testing.T) {
	cache := NewLruCache(1)
	_, ok := cache.Get("invalid")
	if ok {
		t.Error("expected false")
	}
}

func TestSetEmpty(t *testing.T) {
	cache := NewLruCache(1)
	exp := 1000
	result := cache.Set("key", exp)
	if result != nil {
		t.Error("nil expected")
	}
	value,_ := cache.Get("key")
	if value != exp {
		t.Errorf("Expected %v got %v", exp, value)
	}
}

func TestSetTwiceSameKey(t *testing.T) {
	cache := NewLruCache(1)
	cache.Set("same", 1000)
	cache.Set("same", 2000)
	item, ok := cache.Get("same")
	if !ok || item != 2000 {
		t.Error("invalid should be 2000")
	}
}

func TestSetTwoWithOneLimit(t *testing.T) {
	cache := NewLruCache(1)
	cache.Set("first", 1000)
	cache.Set("second", 2000)
	item, _ := cache.Get("first")
	if item != nil {
		t.Error("invalid should be nil got: ")
	}
	item, _ = cache.Get("second")
	if item != 2000 {
		t.Error("invalid should be 2000")
	}
}

func TestDeleteInvalidKey(t *testing.T) {
	cache := NewLruCache(1)
	cache.Del("invalid")
}

func TestSetAndDelete(t *testing.T) {
	cache := NewLruCache(1)
	cache.Set("k", 0)
	cache.Del("k")
	if item,_ := cache.Get("k"); item != nil {
		t.Fail()
	}
}
