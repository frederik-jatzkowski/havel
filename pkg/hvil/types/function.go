package types

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Function struct {
	tool.Node[Function]

	Parameters  tool.List[Type] `parser:"'func' '(' @@ ')'"`
	ReturnValue Type            `parser:"( '->' '(' @@ ')' )?"`
}

func (node *Function) CanBeAssigned(other Type) bool {
	return node.Equals(other)
}

func (node *Function) CanBeAssignedDetailed(other Type) error {
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

	if err := node.ReturnType().CanBeAssignedDetailed(otherFn.ReturnType()); err != nil {
		return fmt.Errorf("function return value mismatch: %w", err)
	}

	return nil
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

	if err := node.ReturnType().EqualsDetailed(otherFn.ReturnType()); err != nil {
		return fmt.Errorf("function return value mismatch: %w", err)
	}

	return nil
}

func (node *Function) Bytes() int {
	return 8
}

func (node *Function) ReturnType() Type {
	if node.ReturnValue == nil {
		return &Void{}
	}

	return node.ReturnValue
}

func (node *Function) CanDoArithmetics() bool {
	return false
}
