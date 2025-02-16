package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterValidate(t *testing.T) {
	type TestStruct struct {
		TestField int `validate:"nonzero"`
	}

	validator := New(ValidatorOptions{})
	validator.RegisterOperator(reflect.Int, "nonzero", func(obj any, _ string) bool {
		return obj.(int) != 0
	})

	err := validator.Register(TestStruct{})
	require.NoError(t, err)

	// Invalid object
	ok, err := validator.Validate(TestStruct{})
	require.NoError(t, err)
	require.False(t, ok)

	// Valid object
	ok, err = validator.Validate(TestStruct{
		TestField: 1,
	})
	require.NoError(t, err)
	require.True(t, ok)
}

func TestNoRegister(t *testing.T) {
	type TestStruct struct {
		TestField bool `validate:"something"`
	}

	validator := New(ValidatorOptions{})
	err := validator.Register(TestStruct{})
	require.Error(t, err)
	require.ErrorContains(t, err, "not registered")
}
