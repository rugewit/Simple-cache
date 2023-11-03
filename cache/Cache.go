package cache

import (
	"sync"
	"time"
)

type Cache struct {
	items sync.Map
	close chan struct{}
}

type item struct {
	data    interface{}
	expires int64
}

func New(cleaningInterval time.Duration) *Cache {
	cache := &Cache{
		close: make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(cleaningInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()

				cache.items.Range(func(key, value interface{}) bool {
					item := value.(item)

					if item.expires > 0 && now > item.expires {
						cache.items.Delete(key)
					}

					return true
				})

			case <-cache.close:
				return
			}
		}
	}()

	return cache
}

func (cache *Cache) Get(key interface{}) (interface{}, bool) {
	obj, exists := cache.items.Load(key)

	if !exists {
		return nil, false
	}

	item := obj.(item)

	if item.expires > 0 && time.Now().UnixNano() > item.expires {
		return nil, false
	}

	return item.data, true
}

func (cache *Cache) Set(key interface{}, value interface{}, duration time.Duration) {
	var expires int64

	if duration > 0 {
		expires = time.Now().Add(duration).UnixNano()
	}

	cache.items.Store(key, item{
		data:    value,
		expires: expires,
	})
}

func (cache *Cache) Range(f func(key, value interface{}) bool) {
	now := time.Now().UnixNano()

	fn := func(key, value interface{}) bool {
		item := value.(item)

		if item.expires > 0 && now > item.expires {
			return true
		}

		return f(key, item.data)
	}

	cache.items.Range(fn)
}

func (cache *Cache) Delete(key interface{}) {
	cache.items.Delete(key)
}

func (cache *Cache) Close() {
	cache.close <- struct{}{}
	cache.items = sync.Map{}
}
