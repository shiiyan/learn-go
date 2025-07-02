package main

import (
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func TestWaitGroup(t *testing.T) {
	t.Run("basic usage", func(t *testing.T) {
		var wg sync.WaitGroup
		results := make([]int, 5)

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				results[index] = index * 2
			}(i)
		}

		wg.Wait()

		// Verify all goroutines completed
		for i, v := range results {
			if v != i*2 {
				t.Errorf("Expected %d, got %d at index %d", i*2, v, i)
			}
		}
	})

	t.Run("wait timeout", func(t *testing.T) {
		var wg sync.WaitGroup
		done := make(chan bool)

		wg.Add(1)

		go func() {
			wg.Wait()
			done <- true
		}()

		// Simulate work
		go func() {
			time.Sleep(100 * time.Millisecond)
			wg.Done()
		}()

		select {
		case <-done:
			// Success
		case <-time.After(200 * time.Millisecond):
			t.Error("WaitGroup didn't complete in time")
		}
	})
}

// GOEXPERIMENT=synctest go test -v
func TestBasicSyncTest(t *testing.T) {
	synctest.Run(func() {
		start := time.Now()

		done := make(chan bool)

		// This goroutine will run in controlled time
		go func() {
			time.Sleep(5 * time.Second)
			done <- true
		}()

		<-done

		// In real time, this would take 5 seconds
		// In synctest, it completes instantly
		elapsed := time.Since(start)
		t.Logf("Elapsed time: %v", elapsed) // Will show 5s even though test ran instantly
	})
}

func TestRaceCondition(t *testing.T) {
	synctest.Run(func() {
		var counter int
		done := make(chan bool, 2)

		// Two goroutines incrementing counter
		go func() {
			temp := counter
			counter = temp + 1
			done <- true
		}()

		go func() {
			temp := counter
			counter = temp + 1
			done <- true
		}()

		<-done
		<-done

		// Without proper synchronization, counter might be 1 or 2
		// synctest makes this deterministic
		if counter != 2 {
			t.Errorf("Race condition detected: counter = %d", counter)
		}
	})
}
