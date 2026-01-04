package types

import "fmt"

type Type interface {
	fmt.Stringer
	CanBeAssigned(Type) bool
	Equals(Type) bool
	EqualsDetailed(Type) error
	Bytes() int
}
