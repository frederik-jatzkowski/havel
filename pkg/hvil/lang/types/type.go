package types

import "fmt"

type Type interface {
	fmt.Stringer
	CanBeAssigned(Type) bool
	Bytes() int
}
