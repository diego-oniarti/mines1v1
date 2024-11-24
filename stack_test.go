package main

import (
	"testing"

	"github.com/diego-oniarti/mines1v1/gamemodes"
)

func TestStack(t *testing.T) {
    s := gamemodes.NewStack[int]()
    s.Push(5)
    s.Push(2)
    if s.Pop() != 2 {
        t.Fatal()
    }
    if s.Pop() != 5 {
        t.Fatal()
    }
    if s.Len() != 0 {
        t.Fatal(s.Len())
    }

    for i:=0; i<100; i++ {
        s.Push(i)
    }
    for i:=99; i>=0; i-- {
        a := s.Pop()
        if a!=i {
            t.Fatal(a,i)
        }
    }
}
