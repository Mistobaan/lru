// Package lru implements a high efficient LRU cache that stores map[string][]byte and with expiration date
package lru

import (
	"sync"
	"time"
)

// Item is the container for the elements inside the LRU cache
type Item struct {
	next *Item
	prev *Item

	key   string
	value []byte

	lastAccess time.Time
	expiresIn  time.Duration
}

// Cache is a typical LRU cache implementation. When an element is
// accessed it is promoted to the head of the list, and when space is
// needed the element at the tail of the list (the least recently used
// element) is evicted.
type Cache struct {
	mu sync.Mutex

	head *Item
	tail *Item

	table map[string]*Item

	size     int64
	capacity int64

	// don't thrash the allocated items
	pool *Item
}

// NewCache creates a cache that will keep only the last `size` element in memory. As mapper it uses a standard map[Key]*item
func NewCache(capacity int64) *Cache {
	return &Cache{
		head:     nil,
		tail:     nil,
		size:     0,
		capacity: capacity,
		pool:     nil,
	}
}

func (c *Cache) pushFront(it *Item) {
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

func (c *Cache) popTail() {
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

func (c *Cache) newItem(key string, value []byte) *Item {
	// check the pool
	if c.pool == nil {
		return &Item{
			key:   key,
			value: value,
			next:  nil,
			prev:  nil,
		}
	}

	r := c.pool
	c.pool = r.next
	r.key = key
	r.value = value
	return r
}

func (c *Cache) delItem(it *Item) {
	a := c.pool
	c.pool = it
	it.next = a
	it.prev = nil
	it.value = nil
	it.expiresIn = 0
	it.key = ""
}

func (c *Cache) pop(it *Item) {
	if it == c.tail {
		c.popTail()
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
func (c *Cache) Set(key string, value []byte) {

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.capacity == 0 {
		return
	}

	item, ok := c.table[key]
	if !ok {
		c.addNew(key, value)
	} else {
		c.updateInPlace(item, key, value)
	}
}

func (c *Cache) addNew(key string, value []byte) {
	item := c.newItem(key, value)
	c.table[key] = item
	c.pushFront(item)
	c.checkCapacity()
}

func (c *Cache) updateInPlace(item *Item, key string, value []byte) {
	valueSize := int64(cap(value))
	sizeDiff := valueSize - int64(cap(item.value))
	item.value = value
	c.size += sizeDiff
	c.pushFront(item)
	c.checkCapacity()
}

func (c *Cache) checkCapacity() {
	for c.size > c.capacity {
		item := c.tail
		c.delete(item.key)
		c.size -= int64(cap(item.value))
	}
}

// Get Gets the latest value of key if available.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.capacity == 0 {
		return nil, false
	}

	item, ok := c.table[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(item.lastAccess.Add(item.expiresIn)) {
		c.delete(key)
		return nil, false
	}

	c.pushFront(item)
	return item.value, true
}

// Delete Deletes a key from the cache. no action is taken if the key is not found .
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	c.delete(key)
	c.mu.Unlock()
}

// internal delete only. Use with locks
func (c *Cache) delete(key string) {
	it, ok := c.table[key]
	if !ok {
		return
	}
	c.pop(it)
	delete(c.table, key)

}

// Size returns the number of items in the cache
func (c *Cache) Size() int64 {
	return c.size
}
