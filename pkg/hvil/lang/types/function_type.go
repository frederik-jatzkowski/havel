package types

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type FunctionType struct {
	tool.Node[FunctionType]

	Parameters  tool.List[Type] `parser:"'func' '(' @@ ')'"`
	ReturnValue Type            `parser:"( '=>' @@ )?"`
}

func (node *FunctionType) CanBeAssigned(other Type) bool {
	return node.Equals(other)
}

func (node *FunctionType) Equals(other Type) bool {
	return node.EqualsDetailed(other) == nil
}

func (node *FunctionType) EqualsDetailed(other Type) error {
	otherFn, ok := other.(*FunctionType)
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

	if err := node.ReturnValue.EqualsDetailed(otherFn.ReturnValue); err != nil {
		return fmt.Errorf("function return value mismatch: %w", err)
	}

	return nil
}

func (node *FunctionType) Bytes() int {
	return 8
}
