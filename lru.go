// Package lru implements an LRU cache.
package lru

import (
	"fmt"
)

// TODO:
// - add a perfect hash function
// - allow to resize the cache

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key interface{}

type Item struct {
	next  *Item
	prev  *Item
	key   Key
	value interface{}
}

type Cache struct {
	table   Mapper
	head     *Item
	tail     *Item
	free     uint
	capacity uint
	item_pool *Item
}

type Mapper interface {
	GetItem(Key) (*Item, bool)
	SetItem(Key, *Item)
	DelItem(Key)
}

// Default Mapper
type DefaultMapper map[Key]*Item

func (m DefaultMapper) GetItem(k Key) (*Item, bool) {
	item, ok := m[k]
	return item, ok
}

func (m DefaultMapper) SetItem(k Key, value *Item) {
	m[k] = value
}

func (m DefaultMapper) DelItem(k Key) {
	delete(m, k)
}

// New creates a cache that will keep only the last `size` element in memory. As mapper it uses a standard map[Key]*item
func New(size uint) *Cache {
	return NewWithMapper(size, &DefaultMapper{})
}

// NewWithMapper creates a cache that will keep only the last `size` element in memory and will use the provided mapper
func NewWithMapper(size uint, mapper Mapper) *Cache {
	return &Cache{
		table:   mapper,
		head:     nil,
		tail:     nil,
		free:     size,
		capacity: size,
		item_pool: nil,
	}
}

func (c *Cache) push_front(it *Item) {
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
	it := tail

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
	c.delItem(it)
}

func (c *Cache) newItem(key, value interface{}) *Item {
	if c.item_pool == nil {
		return &Item{
			key:   key,
			value: value,
			next:  nil,
			prev:  nil,
		}
	} else {
		r := c.item_pool
		c.item_pool = r.next
		r.key = key
		r.value = value
		return r
	}
}

func (c *Cache) delItem(it *Item) {
	a := c.item_pool
	c.item_pool = it
	it.next = a
	it.prev = nil
	it.value = nil
	it.key = nil
}

func (c *Cache) pop(it *Item) {
	if it == c.tail {
		c.pop_tail()
	} else if it == c.head {
		// Head -> it ->  a
		c.head = it.next
		if c.head != nil {
			c.head.prev = nil
		}
		c.delItem(it)
	} else {
		// A -> B -> C: remove B
		A := it.prev
		C := it.next
		A.next = C
		C.prev = A
		c.delItem(it)
	}
}

// Set sets a value in the cache.
func (c *Cache) Set(key, value interface{}) error {
	// this does not work if capacity is zero
	if c.capacity == 0 {
		return fmt.Errorf("Can't set to a zero capacity")
	}

	if c.free < 1 {
		c.table.DelItem(c.tail.key)
		c.pop_tail()
		c.free += 1
	}

	it, ok := c.table.GetItem(key)

	if !ok {
		it = c.newItem(key, value)
		c.table.SetItem(key, it)
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

	item, ok := c.table.GetItem(key)
	if !ok {
		return nil, false
	}

	c.push_front(item)

	return item.value, true
}

// Del Deletes a key from the cache. no action is taken if the key is not found .
func (c *Cache) Del(key interface{}) {
	it, ok := c.table.GetItem(key)
	if !ok {
		return
	}
	c.pop(it)
	c.free += 1
	c.table.DelItem(key)
}

// Len returns the number of items in the cache
func (c *Cache) Len() uint {
	return c.capacity - c.free
}
