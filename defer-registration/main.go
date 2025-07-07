package main

import "fmt"

func testDeferStack() {
	defer fmt.Println("Deferred: 1")
	defer fmt.Println("Deferred: 2")
	defer fmt.Println("Deferred: 3")
	fmt.Println("Executing function")
}

func testDeferEarlyReturn() {
	if true {
		fmt.Println("Early return condition met")
		return
	}

	defer fmt.Println("Deferred: 1")
	fmt.Println("Executing function")
}

func main() {
	testDeferStack()
	testDeferEarlyReturn()
	fmt.Println("Main function completed")
}
