package main

import (
	"fmt"
	"runtime"
)

// docker run --cpus=2 -v $(pwd):/app golang:1.24-alpine go run /app/main.go

func main() {
	maxProcs := runtime.GOMAXPROCS(0) // returns the current value
	fmt.Println("NumCPU:", runtime.NumCPU())
	fmt.Println("GOMAXPROCS:", maxProcs)
}
