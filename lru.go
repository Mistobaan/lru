package lru

// - allow to specify an hash_function

type Cache struct {
}

func NewLruCache() *Cache {
	return &Cache{}
}

func (c *Cache) Get(key string) interface{} {
	return nil
}
