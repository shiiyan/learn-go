package main

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	FirstName     string     `validate:"required"`
	LastName      string     `validate:"required"`
	Age           uint8      `validate:"gte=0,lte=130"`
	Email         string     `validate:"required,email"`
	Gender        string     `validate:"oneof=male female"`
	FavoriteColor string     `validate:"iscolor"`
	Addresses     []*Address `validate:"required,dive,required"`
}

type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
}

var validate *validator.Validate

func main() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validateStruct()
}

func validateStruct() {
	address := &Address{
		Street: "Eavesdown Docks",
	}
	address2 := &Address{
		City: "Tokyo",
	}

	user := &User{
		FirstName:     "Badger",
		LastName:      "Smith",
		Age:           135,
		Gender:        "male",
		Email:         "Badger.Smith.com",
		FavoriteColor: "#000-",
		Addresses:     []*Address{address, address2},
	}

	err := validate.Struct(user)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				fmt.Println(e.Namespace())
				fmt.Println(e.Field())
				fmt.Println(e.StructNamespace())
				fmt.Println(e.StructField())
				fmt.Println(e.Tag())
				fmt.Println(e.ActualTag())
				fmt.Println(e.Kind())
				fmt.Println(e.Type())
				fmt.Println(e.Value())
				fmt.Println(e.Param())
				fmt.Println()
			}
		}
	}
}
