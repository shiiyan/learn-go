package main

import (
	"fmt"
	"sync"
)

// go run -race main.go
func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var x int

	wg.Add(2)
	go func() {
		mu.Lock()
		x = 1
		mu.Unlock()
		wg.Done()
	}()
	go func() {
		mu.Lock()
		fmt.Println(x)
		mu.Unlock()
		wg.Done()
	}()
	wg.Wait()
}
