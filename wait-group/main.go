package main

import (
	"fmt"
	"sync"
)

// func (wg *WaitGroup) Go(f func()) {
//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         f()
//     }()
// }

func main() {
	var wg sync.WaitGroup

	wg.Add(1) // Increment counter
	go func() {
		defer wg.Done()
		fmt.Println("Hello from goroutine!")
	}()

	wg.Add(1) // Increment counter
	go func() {
		defer wg.Done()
		fmt.Println("Hello from another goroutine!")
	}()

	wg.Wait()
	fmt.Println("All goroutines completed!")
}
