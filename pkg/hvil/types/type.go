package types

import "fmt"

type Type interface {
	fmt.Stringer
	Equals(Type) bool
	EqualsDetailed(Type) error
	Bytes() int
	CanDoArithmetics() bool
	Dereference(fields []uint) (t Type, offset uint, err error)
}
