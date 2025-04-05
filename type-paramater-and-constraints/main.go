package main

import "fmt"

type Constraint interface {
	~[]byte | ~string
	Hash() uint64
}

func at[T Constraint](s T, i int) byte {
	return s[i]
}

type MyString string

func (ms MyString) Hash() uint64 {
	var hash uint64
	for i := range len(ms) {
		hash = hash*31 + uint64(ms[i])
	}
	return hash
}


func main() {
	var ms MyString = "Hello"
	fmt.Println("ms at 2:", string(at(ms, 2)))
	fmt.Println("ms hash:", ms.Hash())
}
