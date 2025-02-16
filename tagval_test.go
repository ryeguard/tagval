package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterValidate(t *testing.T) {
	type TestStruct struct {
		TestField int `validate:"nonzero"`
	}

	validator := New(ValidatorOptions{})
	validator.Register("nonzero", TestStruct{}, func(a any) (bool, error) {
		return a.(TestStruct).TestField != 0, nil
	})

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

func TestRegisterValidateWithTagValue(t *testing.T) {
	type TestStruct struct {
		TestField int `validate:"age#18"`
	}

	validator := New(ValidatorOptions{})
	validator.Register("age", TestStruct{}, func(a any) (bool, error) {
		tagVal, err := strconv.Atoi(validator.TagValue(a, "age"))
		if err != nil {
			return false, fmt.Errorf("atoi: %w", err)
		}
		return a.(TestStruct).TestField >= tagVal, nil
	})

	// Invalid object
	ok, err := validator.Validate(TestStruct{TestField: 15})
	require.NoError(t, err)
	require.False(t, ok)

	// Valid object
	ok, err = validator.Validate(TestStruct{
		TestField: 18,
	})
	require.NoError(t, err)
	require.True(t, ok)
}

func TestNoRegister(t *testing.T) {
	type TestStruct struct {
		TestField bool `validate:"something"`
	}

	validator := New(ValidatorOptions{})
	_, err := validator.Validate(TestStruct{
		TestField: true,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "rule not registered")

}

func TestNoTag(t *testing.T) {
	type TestStruct struct {
		OnlyJsonTag bool `json:"only_json_tag"`
	}

	validator := New(ValidatorOptions{})
	validator.Register("some_other_tag", TestStruct{}, func(a any) (bool, error) {
		return a.(TestStruct).OnlyJsonTag, nil
	})

	_, err := validator.Validate(TestStruct{
		OnlyJsonTag: true,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "struct with registered function was not tagged")
}
