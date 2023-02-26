package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int
	head *ListItem
	tail *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
	}

	if l.head != nil {
		newItem.Next = l.head
		l.head.Prev = newItem
		l.head = newItem
	} else {
		l.head = newItem
		l.tail = newItem
	}
	l.len++
	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
	}
	if l.tail != nil {
		newItem.Prev = l.tail
		l.tail.Next = newItem
		l.tail = newItem
	} else {
		l.head = newItem
		l.tail = newItem
	}
	l.len++
	return l.tail
}

func (l *list) Remove(i *ListItem) {
	if l.len == 1 {
		l.head = nil
		l.tail = nil
		l.len--
		return
	}

	if l.head == i {
		l.head = l.head.Next
		l.head.Prev = nil

		i.Next = nil
		i.Prev = nil
		l.len--
		return
	}

	if l.tail == i {
		l.tail = l.tail.Prev
		l.tail.Next = nil

		i.Next = nil
		i.Prev = nil

		l.len--
		return
	}

	prevElement := i.Prev
	nextElement := i.Next

	prevElement.Next = nextElement
	nextElement.Prev = prevElement

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}
	if l.tail == i {
		l.tail = l.tail.Prev
		l.tail.Next = nil

		l.head.Prev = i

		i.Next = l.head
		i.Prev = nil

		l.head = i
		return
	}
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev

	l.head.Prev = i
	i.Next = l.head
	i.Prev = nil

	l.head = i
}

func NewList() List {
	return new(list)
}
