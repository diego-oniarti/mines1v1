package main;

type Stack[T any] struct {
    items []T;
    length uint;
}

func (s *Stack[T]) push(item T) {
    s.items = append(s.items, item);
    s.length++;
}

func (s *Stack[T]) pop() T {
    s.length--;
    return s.items[s.length];
}

func (s *Stack[T]) len() int {
    return len(s.items);
}

func NewStack[T any]() Stack[T] {
    return Stack[T]{};
}
