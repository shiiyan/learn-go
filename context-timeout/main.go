package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

type HTTPFetcher struct {
	client *http.Client
}

func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (f *HTTPFetcher) Fetch(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Step 1: This establishes TCP connection and sends HTTP request
	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}

	// Step 2: At this point:
	// - TCP connection is established
	// - HTTP headers have been received
	// - Response body might still be streaming from server
	defer resp.Body.Close()

	// Step 3: This reads ALL remaining data from the response body
	// - Might require multiple TCP packets
	// - Could be chunked transfer encoding
	// - Reads until server sends EOF
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
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

	fetcher := NewHTTPFetcher(3 * time.Second)
	fetchResult, err := fetcher.Fetch("https://httpbin.org/delay/5")
	if err != nil {
		fmt.Printf("Fetch error: %v\n", err)
	} else {
		fmt.Printf("Fetch success: %s\n", fetchResult)
	}
}
