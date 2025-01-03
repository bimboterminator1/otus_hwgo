package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	lock     sync.Mutex
}

// In order to achieve constant complexity
// during oldest element ousting.
type cacheValue struct {
	value interface{}
	key   Key
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	item, ok := c.items[key]

	if !ok {
		listItem := c.queue.PushFront(cacheValue{value: value, key: key})
		c.items[key] = listItem

		for c.queue.Len() > c.capacity {
			oldCacheValue, _ := c.queue.Back().Value.(cacheValue)
			c.queue.Remove(c.queue.Back())
			delete(c.items, oldCacheValue.key)
		}

		return false
	}

	item.Value = cacheValue{value: value, key: key}
	c.queue.MoveToFront(item)

	return true
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	item, ok := c.items[key]

	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)
	val, _ := item.Value.(cacheValue)

	return val.value, true
}

func (c *lruCache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for k, item := range c.items {
		c.queue.Remove(item)
		delete(c.items, k)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		lock:     sync.Mutex{},
	}
}
