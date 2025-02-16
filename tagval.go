package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Validator interface {
	Validate(any) (bool, error)
	Register(any) error
	RegisterOperator(reflect.Kind, string, OperatorFunc)
}

type ValidatorOptions struct {
	StructTag string
}

type structValidator struct {
	structTag string
	rules     map[string]map[string]ValidationFunc
	operators map[string]map[reflect.Kind]OperatorFunc // map[operator][type]
}

var _ Validator = &structValidator{}

func (v *structValidator) Validate(obj any) (ok bool, err error) {
	val := reflect.ValueOf(obj)
	structName := val.Type().Name()

	for i := range val.NumField() {
		field := val.Type().Field(i)

		if rule, ok := v.rules[structName][field.Name]; ok {
			if valid, err := rule(val.Field(i).Interface()); err != nil {
				return false, err
			} else if !valid {
				fmt.Printf("field %v failed validation\n", field.Name)
				return false, nil
			}
		}
	}
	return true, nil
}

type ValidationFunc func(obj any) (bool, error)
type OperatorFunc func(obj any, val string) bool

func (v *structValidator) Register(obj any) error {
	val := reflect.ValueOf(obj)
	structName := val.Type().Name()

	for i := range val.NumField() {
		field := val.Type().Field(i)

		validationTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		fieldKind := field.Type.Kind()

		for _, tag := range strings.Split(validationTag, ",") {
			opName := strings.Split(tag, ":")[0]

			if _, ok := v.operators[opName]; !ok {
				return fmt.Errorf("operator '%v' not registered", opName)
			}

			tagValue := ""
			if len(strings.Split(tag, ":")) > 1 {
				tagValue = strings.Split(tag, ":")[1]
			}

			op, ok := v.operators[opName][fieldKind]
			if !ok {
				return fmt.Errorf("operator '%v' not registered for kind '%v'", opName, fieldKind)
			}

			if v.rules[structName] == nil {
				v.rules[structName] = map[string]ValidationFunc{}
			}

			v.rules[structName][field.Name] = func(obj any) (bool, error) {
				return op(obj, tagValue), nil
			}
		}
	}

	return nil
}

func (v *structValidator) RegisterOperator(typ reflect.Kind, name string, fn OperatorFunc) {
	if v.operators[name] == nil {
		v.operators[name] = map[reflect.Kind]OperatorFunc{}
	}

	v.operators[name][typ] = fn
}

func New(opts ValidatorOptions) Validator {
	if opts.StructTag == "" {
		opts.StructTag = "validate"
	}

	return &structValidator{
		structTag: opts.StructTag,
		rules:     map[string]map[string]ValidationFunc{},
		operators: map[string]map[reflect.Kind]OperatorFunc{},
	}
}
