package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	// Проверяем не превысили ли размер кэша
	l.Lock()
	defer l.Unlock()
	overrun := l.capacity == l.queue.Len()

	newValue := cacheItem{
		key:   key,
		value: value,
	}

	// Проверяем, если ли ключ в словаре
	item, keyExists := l.items[key]

	// Обновляем уже добавленный
	if keyExists {
		item.Value = newValue
		l.queue.MoveToFront(item)
		return true
	}

	if overrun {
		// Удаляем последний элемент из очереди
		lastElement := l.queue.Back()

		// Удаляем из словаря и очереди
		delete(l.items, lastElement.Value.(cacheItem).key)
		l.queue.Remove(lastElement)

		// Добавляем новый ключ
		l.queue.PushFront(newValue)
		l.items[key] = l.queue.Front()

		return false
	}

	// Добавляем новый ключ
	l.queue.PushFront(newValue)
	l.items[key] = l.queue.Front()
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	// Проверяем, если ли ключ в словаре
	l.Lock()
	defer l.Unlock()
	item, keyExists := l.items[key]

	if keyExists {
		l.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.Lock()
	defer l.Unlock()
	// Создаем новые пустые объекты очереди и мапки
	l.queue = &list{
		len:  0,
		head: nil,
		tail: nil,
	}
	l.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
