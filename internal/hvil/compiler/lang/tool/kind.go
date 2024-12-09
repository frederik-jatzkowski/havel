package tool

import "fmt"

type Kind[T any] struct {
}

func (k Kind[T]) MarshalText() ([]byte, error) {
	var t T
	return []byte(fmt.Sprintf("%T", t)), nil
}
