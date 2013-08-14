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

func TestSetTwiceSameKey(t *testing.T) {
	cache := NewLruCache()
	cache.Set("same", 1000)
	cache.Set("same", 2000)
	if cache.Get("same") != 2000 {
		t.Error("invalid should be 2000")
	}
}

func TestSetTwoWithOneLimit(t *testing.T) {
	cache := NewLruCache()
	cache.Set("first", 1000)
	cache.Set("second", 2000)
	if cache.Get("first") != nil {
		t.Error("invalid should be nil got: ")
	}
	if cache.Get("second") != 2000 {
		t.Error("invalid should be 2000")
	}
}

func TestDeleteInvalidKey(t *testing.T) {
	cache := NewLruCache()
	cache.Del("invalid")
}

func TestSetAndDelete(t *testing.T) {
	cache := NewLruCache()
	cache.Set("k", 0)
	cache.Del("k")
	if cache.Get("k") != nil {
		t.Fail()
	}
}
