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
	len         int
	front, back *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newElement := &ListItem{Value: v}
	if l.front == nil {
		l.front = newElement
		l.back = newElement
	} else {
		newElement.Next = l.front
		l.front.Prev = newElement
		l.front = newElement
	}
	l.len++
	return newElement
}

func (l *list) PushBack(v interface{}) *ListItem {
	newElement := &ListItem{Value: v}
	if l.back == nil {
		l.front = newElement
		l.back = newElement
	} else {
		newElement.Prev = l.back
		l.back.Next = newElement
		l.back = newElement
	}
	l.len++
	return newElement
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
