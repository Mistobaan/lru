package lru

import "testing"

type MyValue []byte

func (mv MyValue) Size() int {
	return cap(mv)
}

/*
func BenchmarkGetVitessLRU(b *testing.B) {
	cache := cache.NewLRUCache(64 * 1024 * 1024)
	value := make(MyValue, 1000)
	cache.Set("stuff", value)
	for i := 0; i < b.N; i++ {
		val, ok := cache.Get("stuff")
		if !ok {
			panic("error")
		}
		_ = val
	}
}
*/

func BenchmarkGet(b *testing.B) {
	cache := NewCache(1024 * 10)
	cache.Set("stuff", []byte("this is a value"))
	for i := 0; i < b.N; i++ {
		val, ok := cache.Get("stuff")
		if !ok {
			panic("error")
		}
		_ = val
	}
}
