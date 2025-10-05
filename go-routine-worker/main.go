package main

import (
	"context"
	"time"
)

func worker(ctx context.Context, out chan<- struct{}) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case out <- struct{}{}:
			default: // avoid blocking the worker if no one is receiving from the channel.
			}
		}
	}
}
