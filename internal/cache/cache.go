package cache

import "sync"

type Cache[T any] interface {
	Set(key string, value T)
	Get(key string) (T, bool)
}

type MapCache[T any] struct {
	mu    sync.RWMutex
	cache map[string]T
}

func NewMapCache[T any](capacity int) *MapCache[T] {
	return &MapCache[T]{
		cache: make(map[string]T, capacity),
	}
}

func (c *MapCache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = value
}

func (c *MapCache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.cache[key]
	return value, ok
}

type tuple[string, T any] struct {
	key   string
	value T
}

type SliceCache[T any] struct {
	mu    sync.RWMutex
	cache []tuple[string, T]
}

func NewSliceCache[T any](capacity int) *SliceCache[T] {
	return &SliceCache[T]{
		cache: make([]tuple[string, T], capacity),
	}
}

func (c *SliceCache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = append(c.cache, tuple[string, T]{key, value})
}
