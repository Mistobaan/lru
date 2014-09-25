package lru

import (
	"bytes"
	"strconv"
	"testing"
)

func TestGetEmpty(t *testing.T) {
	cache := NewCache(1024 * 10)
	_, ok := cache.Get("invalid")
	if ok {
		t.Error("expected false")
	}
}

func TestSetEmpty(t *testing.T) {
	cache := NewCache(1024 * 10)

	exp := []byte{0xFF}
	cache.Set("key", exp)

	value, ok := cache.Get("key")

	if !ok || !bytes.Equal(value, exp) {
		t.Errorf("Expected %v got %v %v", exp, value, cache.table)
	}
}

func TestSetTwiceSameKey(t *testing.T) {
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

/*
func TestSetTwoWithOneLimit(t *testing.T) {
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
	if item != exp2 {
		t.Error("invalid should be 2000")
	}
}

func TestDeleteInvalidKey(t *testing.T) {
	cache := NewCache(1024 * 10)

	cache.Del("invalid")
}

func TestSetAndDelete(t *testing.T) {
	cache := NewCache(1024 * 10)

	cache.Set("k", []byte{})
	cache.Delete("k")
	if item, _ := cache.Get("k"); item != nil {
		t.Fail()
	}
}

type simpleStruct struct {
	int
	string
}

type complexStruct struct {
	int
	simpleStruct
}

var getTests = []struct {
	name       string
	keyToAdd   interface{}
	keyToGet   interface{}
	expectedOk bool
}{
	{"string_hit", "myKey", "myKey", true},
	{"string_miss", "myKey", "nonsense", false},
	{"simple_struct_hit", simpleStruct{1, "two"}, simpleStruct{1, "two"}, true},
	{"simeple_struct_miss", simpleStruct{1, "two"}, simpleStruct{0, "noway"}, false},
	{"complex_struct_hit", complexStruct{1, simpleStruct{2, "three"}},
		complexStruct{1, simpleStruct{2, "three"}}, true},
}

func TestGet(t *testing.T) {
	for _, tt := range getTests {
		cache := NewCache(1024 * 10)

		lru.Set(tt.keyToAdd, 1234)
		val, ok := lru.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestRemove(t *testing.T) {
	cache := NewCache(1024 * 10)

	lru.Set("myKey", 1234)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestRemove failed.  Expected %d, got %v", 1234, val)
	}

	lru.Del("myKey")
	if _, ok := lru.Get("myKey"); ok {
		t.Fatal("TestRemove returned a removed entry")
	}
}
*/

func Benchmark_Insert(b *testing.B) {
	cache := NewCache(1024 * 10)

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		cache.Set(key, []byte{0x0, 0x0, 0x0, 0x0})
	}
}

func Benchmark_Fetch(b *testing.B) {
	cache := NewCache(1024 * 10)

	for i := 0; i < 1000; i++ {
		key := strconv.Itoa(i)
		cache.Set(key, []byte{0x0, 0x0, 0x0, 0x0})
	}

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 1000)
		cache.Get(key)
	}
}
