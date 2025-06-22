package main

import (
	"fmt"
	"time"
)

func loopVariable() {
	values := []int{1, 2, 3, 4, 5}
	for _, val := range values {
		go func() {
			fmt.Printf("%d ", val)
		}()
	}
}

func main() {
	loopVariable()
	time.Sleep(3 * time.Second)
}
