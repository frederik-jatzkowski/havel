package parser

import "fmt"

type Type interface {
	fmt.Stringer
	Equals(Type) bool
}
