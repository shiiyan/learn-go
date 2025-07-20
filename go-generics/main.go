package main

import (
	"cmp"
	"fmt"
	"iter"
)

func Min[T cmp.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Process[T any](item T) T { return item }

func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	var result []T
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	return result
}

type Set[E any] interface {
	Insert(E)
	Delete(E)
	Has(E) bool
	All() iter.Seq[E]
}


func main() {
	fmt.Println("Min of 3 and 5:", Min(3, 5))
	fmt.Println("Min of 7.2 and 2.4:", Min(7.2, 2.4))
	fmt.Println("Min of 'apple' and 'banana':", Min("apple", "banana"))

	numbers := []int{1, 2, 3, 4, 2, 3, 2, 5, 1}
	fmt.Println("Unique numbers:", Unique(numbers))

	strings := []string{"apple", "banana", "apple", "cherry", "banana"}
	fmt.Println("Unique strings:", Unique(strings))

	floats := []float64{1.1, 2.2, 1.1, 3.3, 2.2, 3.3}
	fmt.Println("Unique floats:", Unique(floats))
}
