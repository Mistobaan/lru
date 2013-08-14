package lru

import (
	"testing"
)

func TestGetEmpty(t *testing.T) {
	cache := New(1)
	_, ok := cache.Get("invalid")
	if ok {
		t.Error("expected false")
	}
}

func TestSetEmpty(t *testing.T) {
	cache := New(1)
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
	cache := New(1)
	cache.Set("same", 1000)
	cache.Set("same", 2000)
	item, ok := cache.Get("same")
	if !ok || item != 2000 {
		t.Error("invalid should be 2000")
	}
}

func TestSetTwoWithOneLimit(t *testing.T) {
	cache := New(1)
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
	cache := New(1)
	cache.Del("invalid")
}

func TestSetAndDelete(t *testing.T) {
	cache := New(1)
	cache.Set("k", 0)
	cache.Del("k")
	if item,_ := cache.Get("k"); item != nil {
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
		lru := New(1)
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
	lru := New(1)
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
