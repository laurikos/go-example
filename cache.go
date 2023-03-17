package main

import (
	"sync"
)

type Cache[K comparable, T any] struct {
	data map[K]T
	mu   sync.RWMutex
}

func NewCache[K comparable, T any]() *Cache[K, T] {
	return &Cache[K, T]{
		data: make(map[K]T),
	}
}

func (c *Cache[K, T]) Get(key K) (value T, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok = c.data[key]
	return
}

func (c *Cache[K, T]) Set(key K, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache[K, T]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
