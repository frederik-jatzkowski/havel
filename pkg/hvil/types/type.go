package types

import "fmt"

type Type interface {
	fmt.Stringer
	CanBeAssigned(Type) bool
	CanBeAssignedDetailed(Type) error
	Equals(Type) bool
	EqualsDetailed(Type) error
	Bytes() int
	CanDoArithmetics() bool
}
