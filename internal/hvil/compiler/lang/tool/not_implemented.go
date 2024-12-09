package tool

import "fmt"

type NotImplemented[T any] struct{}

func (n NotImplemented[T]) MarshalText() ([]byte, error) {
	var t T
	return []byte(fmt.Sprintf("NOT IMPLEMENTED: %T", t)), nil
}
