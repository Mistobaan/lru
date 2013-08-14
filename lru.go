package lru

import (
	"fmt"
)

// - allow to specify an hash_function
// - allow to resize the cache
//

type item struct {
	next  *item
	prev  *item
	key   string
	value interface{}
}

type HashFunc func(string) uint

type Cache struct {
	table    []*item
	head     *item
	tail     *item
	hash     HashFunc
	free     uint
	capacity uint
}

func DefaultHashFunc(key string) uint {
	return 0
}

func NewLruCache() *Cache {
	var size uint = 1
	return &Cache{
		table:    make([]*item, size),
		head:     nil,
		tail:     nil,
		hash:     DefaultHashFunc,
		free:     size,
		capacity: size,
	}
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

func (c *Cache) pop_tail() {
	tail := c.tail
	if c.head == c.tail {
		c.head = nil
	}

	if tail != nil {
		// a -> tail
		a := tail.prev
		if a != nil {
			a.next = nil
		}
		c.tail = a
	}
}

func (c *Cache) Set(key string, value interface{}) error {
	idx := c.hash(key)
	// this does not work if capacity is zero
	if c.capacity == 0 {
		return fmt.Errorf("Can't set to a zero capacity")
	}

	if c.free < 1 {
		c.table[c.hash(c.tail.key)] = nil
		c.pop_tail()
		c.free += 1
	}

	c.free -= 1
	it := c.table[idx]

	if nil == it {
		it = &item{
			value: value,
			key:   key,
			next:  nil,
			prev:  nil,
		}
		c.table[idx] = it
		if c.head == nil {
			c.head = it
			c.tail = c.head
		}
	} else {
		c.table[idx].value = value
	}

	c.push_front(it)

	return nil
}

func (c *Cache) Get(key string) interface{} {
	if c.capacity == 0 {
		return nil
	}

	idx := c.hash(key)
	item := c.table[idx]
	if nil == item {
		return nil
	} else if item.key != key {
		// same position but different keys
		return nil
	}

	c.push_front(item)

	return item.value
}

func (c *Cache) Del(key string) {
}
