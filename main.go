package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Validator interface {
	Validate(obj any) (bool, error)
}

type ValidatorOptions struct {
	StructTag string
}

type structValidator struct {
	structTag string
}

var _ Validator = &structValidator{}

func (v *structValidator) Validate(obj any) (bool, error) {
	val := reflect.ValueOf(obj)
	for i := range val.NumField() {
		field := val.Field(i)
		tag := val.Type().Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}

		switch field.Type().Kind() {
		case reflect.Int:

		}

		fmt.Println(field)

		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			switch {
			case rule == "required":
				if field.String() == "" {
					return false, nil
				}
			case strings.HasPrefix(rule, "min="):
				min, err := strconv.Atoi(strings.TrimPrefix(rule, "min="))
				if err != nil {
					return false, err
				}
				if field.CanInt() {
					if field.Int() < int64(min) {
						return false, nil
					}
				}
			}
		}
	}

	return true, nil
}

func New(opts ValidatorOptions) Validator {
	if opts.StructTag == "" {
		opts.StructTag = "validate"
	}
	return &structValidator{
		structTag: opts.StructTag,
	}
}

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
