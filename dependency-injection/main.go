//go:build wireinject
// +build wireinject

package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/wire"
)

type Message string

func NewMessage() Message {
	return Message("Hi there!")
}

type Greeter struct {
	Message Message
	Grumpy  bool
}

var now = time.Now

func NewGreeter(m Message) Greeter {
	var grumpy bool
	if now().Unix()%2 == 0 {
		grumpy = true
	}

	return Greeter{Message: m, Grumpy: grumpy}
}

func (g Greeter) Greet() Message {
	if g.Grumpy {
		return Message("Go away!")
	}

	return g.Message
}

type Event struct {
	Greeter Greeter
}

func NewEvent(g Greeter) (Event, error) {
	if g.Grumpy {
		return Event{}, errors.New("could not create event: event greeter is grumpy")
	}

	return Event{Greeter: g}, nil
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func InitializeEvent() (Event, error) {
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return Event{}, nil
}

func main() {
	e, err := InitializeEvent()
	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(2)
	}

	e.Start()
}
