package types

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Ref struct {
	tool.Node[Ref]

	Ref string `parser:"'ref'"`
}

func (node *Ref) String() string {
	return "ref"
}

func (node *Ref) MarshalText() ([]byte, error) {
	return []byte(node.String()), nil
}

func (node *Ref) Equals(other Type) bool {
	return node.EqualsDetailed(other) == nil
}

func (node *Ref) EqualsDetailed(other Type) error {
	_, ok := other.(*Ref)
	if !ok {
		return fmt.Errorf("%s is not a ref type", other)
	}

	return nil
}

func (node *Ref) Bytes() int {
	return 8
}

func (node *Ref) CanDoArithmetics() bool {
	return false
}

func (node *Ref) Dereference(fields []uint) (Type, uint, error) {
	if len(fields) == 0 {
		return node, 0, nil
	}

	return nil, 0, fmt.Errorf("cannot dereference into %T", node)
}
