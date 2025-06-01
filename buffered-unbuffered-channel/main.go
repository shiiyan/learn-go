package main

import "fmt"

func main() {
	bufferedMsg := make(chan string, 2)

	bufferedMsg <- "buffered"
	bufferedMsg <- "channel"

	fmt.Println(<-bufferedMsg)
	fmt.Println(<-bufferedMsg)

	unbufferedMsg := make(chan string)

	go func() {
		unbufferedMsg <- "unbuffered"
		unbufferedMsg <- "channel"
	}()

	fmt.Println(<-unbufferedMsg)
	fmt.Println(<-unbufferedMsg)
}
