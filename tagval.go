package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Validator interface {
	Validate(obj any) (bool, error)
	Register(string, any, ValidationFunc)
}

type ValidatorOptions struct {
	StructTag string
}

type structField struct {
	structName string
	ruleName   string
}

type structValidator struct {
	structTag       string
	validationFuncs map[structField](ValidationFunc)
}

var _ Validator = &structValidator{}

func (v *structValidator) Validate(obj any) (bool, error) {
	val := reflect.ValueOf(obj)

	tagged := false

	for i := range val.NumField() {
		tag := val.Type().Field(i).Tag.Get(v.structTag)
		if tag == "" {
			continue
		} else {
			tagged = true
		}

		ruleNames := strings.Split(tag, ",")

		checked := map[string]bool{}

		for _, name := range ruleNames {
			rule, ok := v.validationFuncs[structField{
				structName: reflect.TypeOf(obj).Name(),
				ruleName:   name,
			}]
			if !ok {
				return false, fmt.Errorf("rule not registered")
			}

			valid, err := rule(obj)
			if err != nil {
				return false, fmt.Errorf("rule error: %w", err)
			}

			checked[name] = valid
		}
		for _, check := range checked {
			if !check {
				return false, nil
			}
		}
	}

	if !tagged {
		return false, fmt.Errorf("struct with registered function was not tagged")
	}

	return true, nil
}

type ValidationFunc func(obj any) (bool, error)

func (v *structValidator) Register(name string, obj any, fun ValidationFunc) {
	if v.validationFuncs == nil {
		v.validationFuncs = map[structField]ValidationFunc{}
	}
	structName := reflect.TypeOf(obj).Name()
	v.validationFuncs[structField{
		structName: structName,
		ruleName:   name,
	}] = fun
}

func New(opts ValidatorOptions) Validator {
	if opts.StructTag == "" {
		opts.StructTag = "validate"
	}
	return &structValidator{
		structTag: opts.StructTag,
	}
}
