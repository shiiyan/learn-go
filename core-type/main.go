package main

import (
	"fmt"
)

type ByteOrString interface {
	~[]byte | ~string
}

func sliceFirst[T ByteOrString](s T) T {
	return s[1:]
}

func main() {
	fmt.Println(sliceFirst("abc"))
}
