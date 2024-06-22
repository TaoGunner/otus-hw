package hw04lrucache

type List interface {
	Len() int                          // длина списка
	Front() *ListItem                  // первый элемент списка
	Back() *ListItem                   // последний элемент списка
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
}

type ListItem struct {
	Value interface{} // значение
	Next  *ListItem   // следующий элемент
	Prev  *ListItem   // предыдущий элемент
}

type list struct {
	count int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

// PushFront добавляет значение в начало
func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  l.Front(),
	}

	if l.Len() == 0 {
		l.back = item
	} else {
		l.Front().Prev = item
	}
	l.front = item
	l.count++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Prev:  l.Back(),
	}

	if l.Len() == 0 {
		l.front = item
	} else {
		l.Back().Next = item
	}
	l.back = item
	l.count++

	return item
}

// Remove удаляет существующий элемент из списка
func (l *list) Remove(i *ListItem) {
	switch {
	case l.Len() == 1: // Удаление единственного элемента списка
		l.front, l.back = nil, nil
	case i.Prev == nil: // Удаление первого элемента
		l.front = i.Next
		l.front.Prev = nil
	case i.Next == nil: // Удаление последнего элемента
		l.back = i.Prev
		l.back.Next = nil
	default: // Удаление не крайнего элемента
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
