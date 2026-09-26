package main

import (
	"fmt"
)

type Stack[T comparable] struct {
	vals []T
}

func (s *Stack[T]) Push(val T) {
	s.vals = append(s.vals, val)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.vals) == 0 {
		var zero T
		return zero, false
	}
	top := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return top, true
}

func (s Stack[T]) Empty() bool {
	return len(s.vals) == 0
}

func (s Stack[T]) Contains(val T) bool {
	for _, v := range s.vals {
		if v == val {
			return true
		}
	}
	return false
}
func main() {
	var st Stack[int]
	st.Push(2)
	st.Push(3)
	st.Push(4)
	st.Push(5)

	// for !st.Empty() {
	// 	top, _ := st.Pop()
	// 	fmt.Println(top)
	// }

	fmt.Println(st.Contains(3))

}
