package main

import (
	"context"
	"fmt"
	"time"
)

func ReadWithContext(ctx context.Context, in chan int) (int, error) {
	select {
	case v := <-in:
		return v, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func main() {
	ch := make(chan int)

	ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel1()

	result, err := ReadWithContext(ctx1, ch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: got value %d\n", result)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()

	go func() {
		time.Sleep(1 * time.Second)
		ch <- 42
	}()

	result, err = ReadWithContext(ctx2, ch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: got value %d\n", result)
	}
}
