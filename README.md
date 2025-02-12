# tagval
A minimal Go tag struct validator. 

This package is work in progress and subject to change.

## Usage

The intention is for the API to look something like the snippet below. All functionality is not yet implemented. Please see any/all `*_test.go` files for the most up-to-date usage examples.

```go
package your_package

import (
    "fmt"
    "github.com/ryeguard/tagval"
)

type User struct {
    Name string
    Age  int    `validate:"adult"`
}

func main() {
	validator := tagval.New()

	validator.Register("adult", User{}, func(obj any) (bool, error) {
		return obj.(User).Age >= 18, nil
	})

    user := User{
        Name: "John",
        Age:  20,
    }

    ok, err := validator.Validate(user)
    if err != nil {
        panic(err)
    }
    fmt.Println(ok) // true
}
```

