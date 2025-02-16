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
    Name string `validate:"required"`
    Age  int    `validate:">=:18"` 
}

func main() {

    validator := tagval.New()

    validator.RegisterOperator(reflect.String, "required", func(obj any, _ string) bool {
        return obj.(string) != ""
    })

    validator.RegisterOperator(reflect.Int, ">=", func(obj any, value string) bool {
        i, err := strconv.Atoi(value)
        if err != nil {
            return false
        }
        return obj.(int) >= i
	})

    ok, err := validator.Validate(User{
        Name: "John",
        Age:  20,
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(ok) // true
}
```
