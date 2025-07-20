package main

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
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

type HashSet[E comparable] map[E]bool

func (s HashSet[E]) Insert(v E)       { s[v] = true }
func (s HashSet[E]) Delete(v E)       { delete(s, v) }
func (s HashSet[E]) Has(v E) bool     { return s[v] }
func (s HashSet[E]) All() iter.Seq[E] { return maps.Keys(s) }

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

	// Create a new HashSet for integers
	intSet := make(HashSet[int])

	// Insert some values
	intSet.Insert(1)
	intSet.Insert(2)
	intSet.Insert(3)
	intSet.Insert(2) // Duplicate - won't be added twice

	// Check if values exist
	fmt.Println("Has 2:", intSet.Has(2)) // true
	fmt.Println("Has 5:", intSet.Has(5)) // false

	// Print all values using the iterator
	fmt.Print("All values: ")
	for value := range intSet.All() {
		fmt.Print(value, " ")
	}
	fmt.Println()

	// Delete a value
	intSet.Delete(2)
	fmt.Println("After deleting 2, has 2:", intSet.Has(2)) // false

	// Print all values again
	fmt.Print("After deletion: ")
	for value := range intSet.All() {
		fmt.Print(value, " ")
	}
	fmt.Println()

	// Example with strings
	stringSet := make(HashSet[string])
	stringSet.Insert("apple")
	stringSet.Insert("banana")
	stringSet.Insert("apple") // Duplicate

	fmt.Print("String set: ")
	for value := range stringSet.All() {
		fmt.Print(value, " ")
	}
	fmt.Println()
}
