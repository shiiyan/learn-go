package main

import "fmt"

type MyIntFunc func(int) int

func (f MyIntFunc) Double(x int) int {
	return f(x) * 2
}

func addOne(x int) int {
	return x + 1
}

func square(x int) int {
	return x * x
}

func negate(x int) int {
	return -x
}

func main() {
	var addOneFn MyIntFunc = addOne
	fmt.Println(addOneFn.Double(2))

	var squareFn MyIntFunc = square
	fmt.Println(squareFn.Double(2))

	var negateFn MyIntFunc = negate
	fmt.Println(negateFn.Double(2))
}
