package lru

import (
	"fmt"
)

// - allow to specify an hash_function
// - allow to resize the cache

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key interface{}

type item struct {
	next  *item
	prev  *item
	key   Key
	value interface{}
}

type Cache struct {
	table    map[Key]*item
	head     *item
	tail     *item
	free     uint
	capacity uint
}

// NewLruCache creates a cache that will keep only the last `size` element in memory
func NewLruCache(size uint) *Cache {
	return &Cache{
		table:    make(map[Key]*item),
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

func (c *Cache) Set(key, value interface{}) error {
	// this does not work if capacity is zero
	if c.capacity == 0 {
		return fmt.Errorf("Can't set to a zero capacity")
	}

	if c.free < 1 {
		delete(c.table, c.tail.key)
		c.pop_tail()
		c.free += 1
	}

	it, ok := c.table[key]

	if !ok {
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
		it.value = value
	}

	c.push_front(it)

	return nil
}

// Get Gets the latest value of key if available.
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	if c.capacity == 0 {
		return nil, false
	}

	item, ok := c.table[key]
	if !ok {
		return nil, false
	}

	c.push_front(item)

	return item.value, true
}

// Del Deletes a key from the cache. no action is taken if the key is not found .
func (c *Cache) Del(key interface{}) {
	it, ok := c.table[key]
	if !ok {
		return
	}
	c.pop(it)
	c.free += 1
	delete(c.table, key)
}

// Len returns the number of items in the cache
func (c *Cache) Len() uint {
	return c.capacity - c.free
}
