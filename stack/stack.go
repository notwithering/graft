package stack

type Stack[T any] []T

func (s *Stack[T]) Push(item T) {
	*s = append(*s, item)
}

func (s *Stack[T]) Pop() T {
	if len(*s) == 0 {
		var zero T
		return zero
	}

	item := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return item
}

func (s *Stack[T]) Len() int {
	return len(*s)
}

func (s *Stack[T]) Top() T {
	if len(*s) == 0 {
		var zero T
		return zero
	}
	return (*s)[len(*s)-1]
}
