package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу.
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу.
	Clear()                              // Очистить кэш.
}

type lruCache struct {
	capacity int
	queue    List
	mx       *sync.Mutex
	items    map[Key]*ListItem
}

type Item struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		mx:       &sync.Mutex{},
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	if lru.capacity == 0 {
		return false
	}

	item := Item{
		key:   key,
		value: value,
	}

	lru.mx.Lock()
	_, ok := lru.items[key]
	if ok {
		lru.queue.MoveToFront(lru.items[key])
		lru.queue.Front().Value = item
		lru.items[key] = lru.queue.Front()
	} else {
		if lru.queue.Len() >= lru.capacity {
			delete(lru.items, lru.queue.Back().Value.(Item).key)
			lru.queue.Remove(lru.queue.Back())
		}
		lru.items[key] = lru.queue.PushFront(item)
	}
	lru.mx.Unlock()

	return ok
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mx.Lock()
	defer lru.mx.Unlock()

	item, ok := lru.items[key]
	if ok {
		lru.queue.MoveToFront(lru.items[key])
		return item.Value.(Item).value, true
	}

	return nil, false
}

func (lru *lruCache) Clear() {
	lru.mx.Lock()
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
	lru.mx.Unlock()
}
