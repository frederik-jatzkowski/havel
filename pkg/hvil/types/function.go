package types

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Function struct {
	tool.Node[Function]

	Parameters tool.List[Type] `parser:"'func' '(' @@ ')'"`
}

func (node *Function) Equals(other Type) bool {
	return node.EqualsDetailed(other) == nil
}

func (node *Function) EqualsDetailed(other Type) error {
	otherFn, ok := other.(*Function)
	if !ok {
		return fmt.Errorf("%s is not a function type", other)
	}

	if len(otherFn.Parameters.Items) != len(node.Parameters.Items) {
		return fmt.Errorf("function parameter count mismatch: expected %d, got %d", len(node.Parameters.Items), len(otherFn.Parameters.Items))
	}

	for i, param := range node.Parameters.Items {
		if err := param.EqualsDetailed(otherFn.Parameters.Items[i]); err != nil {
			return fmt.Errorf("function parameter %d mismatch: expected %s, got %s", i, param, otherFn.Parameters.Items[i])
		}
	}

	return nil
}

func (node *Function) Bytes() int {
	return 8
}

func (node *Function) CanDoArithmetics() bool {
	return false
}

func (node *Function) Dereference(fields []uint) (Type, uint, error) {
	if len(fields) == 0 {
		return node, 0, nil
	}

	return nil, 0, fmt.Errorf("cannot dereference into %T", node)
}
