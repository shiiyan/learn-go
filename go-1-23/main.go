package main

import (
	"fmt"
	"iter"
	"sync"
)

func rangeIteration() {
	var m sync.Map

	m.Store("alice", 11)
	m.Store("bob", 12)
	m.Store("cindy", 13)

	fmt.Println("go 1.22")
	m.Range(func(key, value any) bool {
		if key == "bob" {
			return false
		}

		fmt.Println(key, value)
		return true
	})

	fmt.Println("go 1.23")
	for key, val := range m.Range {
		if key == "bob" {
			break
		}
		fmt.Println(key, val)
	}
}

func Reversed[V any](s []V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(s[i]) {
				return
			}
		}
	}
}

func PrintAll[V any](s iter.Seq[V]) {
	for v := range s {
		fmt.Print(v, " ")
	}
	fmt.Println()
}

func main() {
	// rangeIteration()

	// i := []int{1, 2, 3, 4, 5}
	// PrintAll(Reversed(i))
	// s := []string{"a", "b", "c"}
	// PrintAll(Reversed(s))

	s := []string{"a", "b", "c"}
	next, stop := iter.Pull(Reversed(s))
	defer stop()

	for {
		v, ok := next()
		if !ok {
			break
		}
		fmt.Print(v, " ")
	}
}
