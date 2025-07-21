package main

import (
	"fmt"
	"log"
)

// +govalid:required
type Person struct {
	// +govalid:required
	Name string `json:"name"`
	// +govalid:email
	Email string `json:"email"`
}

func main() {
	p := &Person{Name: "", Email: "invalid-email"}

	if err := ValidatePerson(p); err != nil {
		log.Printf("Validation error: %v", err)
		return
	}

	fmt.Printf("Person: %+v", p)
}
