package main

import "fmt"

type Greeter interface {
	Greet(name string) string
}

type greeterImpl struct {
	prefix string
}

func (g *greeterImpl) Greet(name string) string {
	return fmt.Sprintf("%s, %s!", g.prefix, name)
}

func NewGreeter(prefix string) Greeter {
	return &greeterImpl{prefix: prefix}
}

type Handler struct {
	Greeter Greeter
}

func NewHandler(g Greeter) *Handler {
	return &Handler{Greeter: g}
}

func (h *Handler) SayHello(name string) {
	fmt.Println(h.Greeter.Greet(name))
}

func main() {
	greeter := NewGreeter("Hello")
	handler := NewHandler(greeter)

	handler.SayHello("World")
}
