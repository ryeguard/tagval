package main

import (
	"fmt"
	"reflect"
	"strconv"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Person struct {
	FirstName string `validate:"required"`
	Title     string `validate:"titleCase"`
	Handle    string `validate:"snakeCase"`
	Age       int    `validate:">=:0,<=:18"`
}

func main() {
	validator := New(ValidatorOptions{})

	validator.RegisterOperator(reflect.String, "required", func(obj any, tagValue string) bool {
		return obj.(string) != ""
	})

	validator.RegisterOperator(reflect.String, "titleCase", func(obj any, tagValue string) bool {
		return obj.(string) == cases.Title(language.English).String(obj.(string))
	})

	validator.RegisterOperator(reflect.Int, ">=", func(obj any, tagValue string) bool {
		val, err := strconv.Atoi(tagValue)
		if err != nil {
			return false
		}
		return obj.(int) >= val
	})

	validator.RegisterOperator(reflect.Int, "<=", func(obj any, tagValue string) bool {
		val, err := strconv.Atoi(tagValue)
		if err != nil {
			return false
		}
		return obj.(int) <= val
	})

	validator.RegisterOperator(reflect.String, "snakeCase", func(obj any, tagValue string) bool {
		for _, r := range obj.(string) {
			if r != '_' && !('a' <= r && r <= 'z') && !('0' <= r && r <= '9') {
				return false
			}
		}
		return true
	})

	err := validator.Register(Person{})
	if err != nil {
		panic(err)
	}

	valid, err := validator.Validate(Person{
		FirstName: "John",
		Title:     "Senior Developer",
		Handle:    "john_doe",
		Age:       18,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("valid:", valid)
}
