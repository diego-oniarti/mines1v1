package gamemodes

import "testing"

func TestStack(t *testing.T) {
    s := NewStack[int]()
    s.Push(5)
    s.Push(2)
    a := s.Pop()
    if a != 2 {
        t.Fatal()
    }
}
