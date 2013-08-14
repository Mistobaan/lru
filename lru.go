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

func (c *Cache) push_front(it *item) {
	// assuming item != null
	if it == c.head {
		// item already in front
		return
	} else if it == c.tail {
		c.tail = it.prev
		c.tail.next = nil
	} else {
		// A -> B -> C: remove B
		A := it.prev
		C := it.next
		A.next = C
		C.prev = A
	}
	it.next = c.head
	it.prev = nil
	c.head = it
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
		if c.head == nil {
			c.head = c.table[idx]
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

	c.push_front(item)

	return item.value
}
