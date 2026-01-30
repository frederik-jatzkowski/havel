package alu

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
)

func resolveBinOpTypesWithTarget[T Operation](node T, left, right instruction.MemoryRead, target types.Type) error {
	if err := resolveBinOpTypes(node, left, right); err != nil {
		return err
	}

	leftType := left.Type()
	rightType := right.Type()

	if !(target.Bytes() == leftType.Bytes() || target.Bytes() == rightType.Bytes()) {
		if leftType.Equals(rightType) {
			return node.Errorf("cannot assign %s to %s", leftType, target)
		}

		return node.Errorf("cannot assign %s or %s to %s", leftType, rightType, target)
	}

	return nil
}

func resolveBinOpTypes[T Operation](node T, left, right instruction.MemoryRead) error {
	leftType := left.Type()
	rightType := right.Type()

	if !leftType.CanDoArithmetics() {
		return left.Errorf("cannot perform alu operations on %s", leftType)
	}

	if !rightType.CanDoArithmetics() {
		return right.Errorf("cannot perform alu operations on %s", rightType)
	}

	if leftType.Bytes() != rightType.Bytes() {
		return node.Errorf("unequally sized parameters %s and %s", leftType, rightType)
	}

	return nil
}
