package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	ctx := t.Context()
	const workers = 4

	out := make(chan struct{}, 1)

	var wg sync.WaitGroup
	for w := 1; w<= workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, out)
		}()
	}

	select {
	case <- out:
		fmt.Println("received")
	case <-time.After(1 * time.Second):
		t.Fatal("no output")
	}

	t.Cleanup(wg.Wait)
}