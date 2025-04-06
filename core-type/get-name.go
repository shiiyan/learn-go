//go:build getname

package main

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Name string
	ID   int
}

type HasName interface {
	Person | Employee
}

func getName[T HasName](v T) string {
	return v.Name //  v.Name undefined (type T has no field or method Name)
}
