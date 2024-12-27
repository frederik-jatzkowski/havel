package tool

import (
	"fmt"
	"reflect"
)

type NotImplemented[T any] struct{}

func (n NotImplemented[T]) MarshalText() ([]byte, error) {
	v := reflect.ValueOf(*new(T))
	if !v.IsValid() {
		return nil, nil
	}

	typeName := v.Type().String()

	return []byte(fmt.Sprintf("NOT IMPLEMENTED: %s", typeName)), nil
}
