package types

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type RefType struct {
	tool.Node[RefType]

	Ref string `parser:"'ref'"`
}

func (node *RefType) String() string {
	return "ref"
}

func (node *RefType) MarshalText() ([]byte, error) {
	return []byte(node.String()), nil
}

func (node *RefType) CanBeAssigned(other Type) bool {
	return node.Equals(other)
}

func (node *RefType) CanBeAssignedDetailed(other Type) error {
	return node.EqualsDetailed(other)
}

func (node *RefType) Equals(other Type) bool {
	return node.EqualsDetailed(other) == nil
}

func (node *RefType) EqualsDetailed(other Type) error {
	_, ok := other.(*RefType)
	if !ok {
		return fmt.Errorf("%s is not a ref type", other)
	}

	return nil
}

func (node *RefType) Bytes() int {
	return 8
}
