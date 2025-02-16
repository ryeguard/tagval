package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Validator interface {
	Validate(obj any) (bool, error)
	Register(string, any, ValidationFunc)
	TagValue(obj any, name string) string
}

type ValidatorOptions struct {
	StructTag      string
	FieldSeparator string
}

type structField struct {
	structName string
	ruleName   string
}

type structValidator struct {
	structTag       string
	fieldSeparator  string
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
			name = strings.Split(name, v.fieldSeparator)[0]
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

func (v *structValidator) TagValue(obj any, name string) string {
	val := reflect.ValueOf(obj)
	for i := range val.NumField() {
		tag := val.Type().Field(i).Tag.Get(v.structTag)

		ruleNames := strings.Split(tag, ",")

		for _, rule := range ruleNames {
			if strings.HasPrefix(rule, name) {
				return strings.Split(rule, v.fieldSeparator)[1]
			}
		}
	}
	return ""
}

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
	if opts.FieldSeparator == "" {
		opts.FieldSeparator = "#"
	}
	return &structValidator{
		structTag:      opts.StructTag,
		fieldSeparator: opts.FieldSeparator,
	}
}
