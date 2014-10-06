package lru

import (
	"bytes"
	"testing"
)

func TestGetEmpty(t *testing.T) {
	t.Parallel()
	cache := NewCache(1024 * 10)
	_, ok := cache.Get("invalid")
	if ok {
		t.Error("expected false")
	}
}

func TestSetEmpty(t *testing.T) {
	t.Parallel()
	cache := NewCache(1024 * 10)

	exp := []byte{0xFF}
	cache.Set("key", exp)

	value, ok := cache.Get("key")

	if !ok || !bytes.Equal(value, exp) {
		t.Errorf("Expected %v got %v %v", exp, value)
	}
}

func TestSetTwiceSameKey(t *testing.T) {
	t.Parallel()

	cache := NewCache(1024 * 10)
	exp := []byte{0xFF}
	exp2 := []byte{0xFE}

	cache.Set("same", exp)
	cache.Set("same", exp2)

	item, ok := cache.Get("same")

	if !ok || !bytes.Equal(item, exp2) {
		t.Error("invalid should be ", exp2)
	}
}

func TestSetTwoWithOneLimit(t *testing.T) {
	t.Parallel()
	cache := NewCache(1024 * 10)

	exp := []byte{0xFF}
	exp2 := []byte{0xFE}

	cache.Set("first", exp)
	cache.Set("second", exp2)
	item, _ := cache.Get("first")
	if item != nil {
		t.Error("invalid should be nil got: ")
	}
	item, _ = cache.Get("second")
	if !bytes.Equal(item, exp2) {
		t.Error("invalid should be 2000")
	}
}

func TestDeleteInvalidKey(t *testing.T) {
	cache := NewCache(1024 * 10)

	cache.Delete("invalid")
}

func TestSetAndDelete(t *testing.T) {
	cache := NewCache(1024 * 10)

	cache.Set("k", []byte{})
	cache.Delete("k")
	if item, _ := cache.Get("k"); item != nil {
		t.Fail()
	}
}
