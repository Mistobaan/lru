package lru

// - allow to specify an hash_function
// - allow to resize the cache
// 


type item struct {
	next  *item
	prev  *item
	value interface{}
}

type HashFunc func(string) uint

type Cache struct {
	table []*item
	head  *item
	tail  *item
	hash  HashFunc
}

func DefaultHashFunc(key string) uint {
	return 0
}

func push_front(item *item){

}


func NewLruCache() *Cache {
	return &Cache{
		table: make([]*item, 1),
		head:  nil,
		tail:  nil,
		hash:  DefaultHashFunc,
	}
}

func (c *Cache) Set(key string, value interface{}) error {
	idx := c.hash(key)
	if c.table[idx] == nil {
		c.table[idx] = &item{
			value: value,
		}
	}
	return nil
}

func (c *Cache) Get(key string) interface{} {
	idx := c.hash(key)
	item := c.table[idx]
	if nil == item {
		return nil
	}

	push_front(item)

	return item.value
}
