package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	FirstName string `validate:"required,titleCase"`
	Email     string `validate:"required,email"`
	Age       int    `validate:"min=0"`
}

func main() {
	person := Person{
		FirstName: "John",
		Email:     "john@example.com",
	}

	t := reflect.TypeOf(person)
	fmt.Printf("%v %v\n", t.Name(), t.Kind())

	for i := range t.NumField() {
		field := t.Field(i)

		tag := field.Tag.Get("validate")
		fmt.Printf("\t%v %v %v %v %v\n", i, field.Name, field.Type.Kind(), field.Type, tag)
	}

	simon := Person{
		FirstName: "Simon",
		Email:     "simon@example.com",
		Age:       0,
	}

	validator := New(ValidatorOptions{})

	ok, err := validator.Validate(simon)
	if err != nil {
		panic(err)
	}

	fmt.Println(ok)
}
