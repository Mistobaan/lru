package lru

// - allow to specify an hash_function

type Item struct {
	next  *Item
	prev  *Item
	Value interface{}
}

type HashFunc func (string) uint 

type Cache struct {
	table []interface{}
	head *Item
	tail *Item
	hash HashFunc
}

func DefaultHashFunc (key string) uint {
	return 0
}


func NewLruCache() *Cache {
	return &Cache{
		table : make([]interface{}, 1),
		head: nil,
		tail: nil,
		hash: DefaultHashFunc,
	}
}

func (c *Cache) Set(key string, value interface{}) error {
	idx := c.hash(key)
	c.table[idx] = value
	return nil
}

func (c *Cache) Get(key string) interface{} {
	idx := c.hash(key)
	return c.table[idx]
}

