package gamemodes

type Stack[T any] struct {
    items  []T
    length uint
}

func (s *Stack[T]) Push(item T) {
    if int(s.length)==len(s.items) {
        s.items = append(s.items, item)
    }else{
        s.items[s.length]=item
    }
    s.length++
}

func (s *Stack[T]) Pop() T {
    s.length--
    return s.items[s.length]
}

func (s *Stack[T]) Len() int {
    return int(s.length)
}

func NewStack[T any]() Stack[T] {
    return Stack[T]{
    	items:  []T{},
    	length: 0,
    }
}
