package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty List", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("remove test", func(t *testing.T) {
		l := NewList()
		a := l.PushFront("a")
		b := l.PushBack("b")
		c := l.PushBack("c")

		l.Remove(a)
		l.Remove(b)

		front := l.Front()
		back := l.Back()

		require.Equal(t, c, front)
		require.Equal(t, c, back)
	})

	t.Run("remove back test", func(t *testing.T) {
		l := NewList()
		a := l.PushFront("a")
		l.PushBack("b")

		l.Remove(l.Back())

		front := l.Front()
		back := l.Back()

		require.Equal(t, a, front)
		require.Equal(t, a, back)
	})

	t.Run("remove front test", func(t *testing.T) {
		l := NewList()
		l.PushFront("a")
		b := l.PushBack("b")

		l.Remove(l.Front())

		front := l.Front()
		back := l.Back()

		require.Equal(t, b, front)
		require.Equal(t, b, back)
	})

}
