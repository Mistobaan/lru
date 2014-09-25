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

	defaultExpiration time.Duration

	// don't thrash the allocated items
	pool *Item
}

// NewCache creates a cache that will keep only the last `size` element in memory. As mapper it uses a standard map[Key]*item
// The default expiration  is 60 second
func NewCache(capacity int64) *Cache {
	return &Cache{
		head:              nil,
		tail:              nil,
		table:             make(map[string]*Item),
		size:              0,
		capacity:          capacity,
		defaultExpiration: 60 * time.Second,
		pool:              nil,
	}
}

// Set sets a value in the cache.
func (c *Cache) Set(key string, value []byte) {
	c.SetExpire(key, value, c.defaultExpiration)
}

// SetExpire sets a key and when it would expire
func (c *Cache) SetExpire(key string, value []byte, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.table[key]
	if !ok {
		c.addNew(key, value, expiration)
	} else {
		c.updateInPlace(item, key, value, expiration)
	}
}

// Get Gets the latest value of key if available.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.table[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(item.lastAccess.Add(item.expiresIn)) {
		c.deleteKey(key)
		return nil, false
	}

	item.lastAccess = time.Now()

	c.pushFront(item)
	return item.value, true
}

func (c *Cache) pushFront(item *Item) {
	// assuming item != null
	if item == c.head {
		// item already in front
		return
	}

	c.pop(item)

	item.next = c.head
	c.head.prev = item
	item.prev = nil
	c.head = item
}

// pop removes the item from the list
func (c *Cache) pop(it *Item) {
	// A -> it -> C: remove B
	if it == c.head {
		c.head = it.next
	}
	if it == c.tail {
		c.tail = it.prev
	}

	A := it.prev
	C := it.next
	if A != nil {
		A.next = C
	}
	if C != nil {
		C.prev = A
	}
	it.next = nil
	it.prev = nil
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
	r.next = nil
	r.prev = nil
	return r
}

func (c *Cache) deleteItem(item *Item) {
	a := c.pool
	c.pool = item
	item.next = a
	item.prev = nil
	item.value = nil
	item.expiresIn = 0
	item.key = ""
}

// internal delete only. Use with locks
func (c *Cache) deleteKey(key string) {
	item, ok := c.table[key]
	if !ok {
		return
	}
	c.pop(item)
	c.deleteItem(item)
	delete(c.table, key)
}

func (c *Cache) addNew(key string, value []byte, expiration time.Duration) {
	item := c.newItem(key, value)

	if c.head == nil {
		c.head = item
		c.tail = item
	} else {
		item.next = c.head
		c.head = item
	}

	item.expiresIn = expiration
	item.lastAccess = time.Now()

	c.table[key] = item

	c.checkCapacity()
}

func (c *Cache) updateInPlace(item *Item, key string, value []byte, expiration time.Duration) {
	valueSize := int64(cap(value))
	sizeDiff := valueSize - int64(cap(item.value))
	item.value = value

	item.lastAccess = time.Now()
	item.expiresIn = expiration

	c.size += sizeDiff
	c.pushFront(item)
	c.checkCapacity()
}

func (c *Cache) checkCapacity() {
	for c.size > c.capacity {
		item := c.tail
		c.deleteKey(item.key)
		c.size -= int64(cap(item.value))
	}
}

// Delete Deletes a key from the cache. no action is taken if the key is not found .
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	c.deleteKey(key)
	c.mu.Unlock()
}

// Size returns the number of items in the cache
func (c *Cache) Size() int64 {
	return c.size
}
