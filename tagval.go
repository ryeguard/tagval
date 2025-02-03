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
