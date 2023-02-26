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

func (l *lruCache) Set(key Key, value interface{}) bool {
	// Проверяем не превысили ли размер кэша
	l.Lock()
	defer l.Unlock()
	overrun := l.capacity == l.queue.Len()

	// Проверяем, если ли ключ в словаре
	item, keyExists := l.items[key]

	// Обновляем уже добавленный
	if keyExists {
		item.Value = value
		l.queue.MoveToFront(item)
		return true
	}

	if overrun {

		// Удаляем последний элемент из очереди
		lastElement := l.queue.Back()

		// Находим под каким ключом хранится значение
		keyInMap := findKeyByValue(l.items, lastElement.Value)

		// Удаляем из словаря и очереди

		delete(l.items, keyInMap)
		l.queue.Remove(lastElement)

		// Добавляем новый ключ
		l.queue.PushFront(value)
		l.items[key] = l.queue.Front()

		return false
	}

	// Добавляем новый ключ
	l.queue.PushFront(value)
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
		return item.Value, true
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

func findKeyByValue(m map[Key]*ListItem, value interface{}) Key {
	for k, v := range m {
		if v.Value == value {
			return k
		}
	}
	return ""
}
