package types

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Void struct {
	tool.Node[Void]
}

func (node *Void) String() string {
	return "void"
}

func (node *Void) Equals(other Type) bool {
	return node.EqualsDetailed(other) == nil
}

func (node *Void) EqualsDetailed(other Type) error {
	_, ok := other.(*Void)
	if !ok {
		return fmt.Errorf("%s is not void", other)
	}

	return nil
}

func (node *Void) Bytes() int {
	return 0
}

func (node *Void) CanDoArithmetics() bool {
	return false
}

func (node *Void) Dereference(fields []uint) (Type, uint, error) {
	if len(fields) == 0 {
		return node, 0, nil
	}

	return nil, 0, fmt.Errorf("cannot dereference into %T", node)
}
