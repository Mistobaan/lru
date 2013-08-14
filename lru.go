package lru

import (
	"fmt"
)

// - allow to specify an hash_function
// - allow to resize the cache

type item struct {
	next  *item
	prev  *item
	key   string
	value interface{}
}

type Cache struct {
	table    map[string]*item
	head     *item
	tail     *item
	free     uint
	capacity uint
}

// NewLruCache creates a cache that will keep only the last `size` element in memory
func NewLruCache(size uint) *Cache {
	return &Cache{
		table:    make(map[string]*item),
		head:     nil,
		tail:     nil,
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
	c.head.prev = it
	it.prev = nil
	c.head = it
}

func (c *Cache) pop_tail() {
	tail := c.tail

	// tail == c.head then is the only item
	if c.head == tail {
		c.head = nil
		c.tail = nil
	}

	// there is at least one element
	// head -> a
	//         |
	// tail -> b
	if tail != nil {
		// a -> tail
		a := tail.prev
		if a != nil {
			a.next = nil
		}
		c.tail = a
	}
}

func (c *Cache) pop(it *item) {
	if it == c.tail {
		c.pop_tail()
	} else if it == c.head {
		// Head -> it ->  a
		c.head = it.next
		if c.head != nil {
			c.head.prev = nil
		}
	} else {
		// A -> B -> C: remove B
		A := it.prev
		C := it.next
		A.next = C
		C.prev = A
	}
}

func (c *Cache) Set(key string, value interface{}) error {
	// this does not work if capacity is zero
	if c.capacity == 0 {
		return fmt.Errorf("Can't set to a zero capacity")
	}

	if c.free < 1 {
		c.table[c.tail.key] = nil
		c.pop_tail()
		c.free += 1
	}

	it := c.table[key]

	if nil == it {
		it = &item{
			value: value,
			key:   key,
			next:  nil,
			prev:  nil,
		}
		c.table[key] = it
		if c.head == nil {
			c.head = it
			c.tail = c.head
		}
		c.free -= 1
	} else {
		c.table[key].value = value
	}

	c.push_front(it)

	return nil
}

// Get Gets the latest value of key if available. Otherwise it returns nil
func (c *Cache) Get(key string) interface{} {
	if c.capacity == 0 {
		return nil
	}

	item := c.table[key]
	if nil == item {
		return nil
	} else if item.key != key {
		// same position but different keys
		return nil
	}

	c.push_front(item)

	return item.value
}

// Del Deletes a key from the cache. no action is taken if the key is not found .
func (c *Cache) Del(key string) {
	it := c.table[key]
	if it == nil {
		return
	}
	c.pop(it)
	c.free += 1
	c.table[key] = nil
}
