package types

import "fmt"

type Type interface {
	fmt.Stringer
	Equals(Type) bool
	EqualsDetailed(Type) error
	Bytes() int
	CanDoArithmetics() bool
}
