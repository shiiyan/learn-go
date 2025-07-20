package main

import (
	"cmp"
	"fmt"
)

func Min[T cmp.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func main() {
	fmt.Println("Min of 3 and 5:", Min(3, 5))
	fmt.Println("Min of 7.2 and 2.4:", Min(7.2, 2.4))
	fmt.Println("Min of 'apple' and 'banana':", Min("apple", "banana"))
}
