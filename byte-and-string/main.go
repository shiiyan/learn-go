package main

import "fmt"

type Bytestring interface {
	~string | ~[]byte
}

func firstChar[T Bytestring](s T) byte {
	return s[0]
}

func main() {
	s := "hello"
	b := []byte("hello")

	fmt.Println("s:", s)
	fmt.Println("b:", b)

	b[0] = 'w'

	fmt.Println("b as string:", string(b))
	fmt.Printf("first character of s: %c\n", firstChar(s))
	fmt.Printf("first character of b: %c\n", firstChar(b))
}
