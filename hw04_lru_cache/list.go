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
	size int
	head *ListItem
	tail *ListItem
}

func (l list) Len() int {
	return l.size
}

func (l list) Front() *ListItem {
	return l.head
}

func (l list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := new(ListItem)
	newItem.Value = v
	newItem.Next = l.head
	if l.head != nil {
		l.head.Prev = newItem
	} else {
		l.tail = newItem
	}
	l.head = newItem
	l.size++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := new(ListItem)
	newItem.Value = v
	newItem.Prev = l.tail
	if l.tail != nil {
		l.tail.Next = newItem
	} else {
		l.head = newItem
	}
	l.tail = newItem
	l.size++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	l.size--
	if l.size == 0 {
		l.head = nil
		l.tail = nil
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.head {
		return
	}

	if i == l.tail {
		l.tail = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i.Next = l.head
	i.Prev = nil

	if l.head != nil {
		l.head.Prev = i
	}

	l.head = i
}

func NewList() List {
	return new(list)
}
