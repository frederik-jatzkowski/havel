package types

import (
	"fmt"
	"strconv"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type ScalarType struct {
	tool.Node[ScalarType]

	Size int `parser:"@Size? 'byte'"`
}

func (node *ScalarType) String() string {
	return strconv.FormatUint(uint64(node.Size), 10) + " byte"
}

func (node *ScalarType) MarshalText() ([]byte, error) {
	return []byte(node.String()), nil
}

func (node *ScalarType) CanBeAssigned(other Type) bool {
	return node.Equals(other)
}

func (node *ScalarType) CanBeAssignedDetailed(other Type) error {
	return node.EqualsDetailed(other)
}

func (node *ScalarType) Equals(other Type) bool {
	return node.EqualsDetailed(other) == nil
}

func (node *ScalarType) EqualsDetailed(other Type) error {
	otherScalar, ok := other.(*ScalarType)
	if !ok {
		return fmt.Errorf("%s is not a scalar type", other)
	}

	if node.Bytes() != otherScalar.Bytes() {
		return fmt.Errorf("expected %d bytes, got %d", node.Bytes(), otherScalar.Bytes())
	}

	return nil
}

func (node *ScalarType) Bytes() int {
	if node.Size == 0 {
		return 1
	}

	return node.Size
}

func (node *ScalarType) CanDoArithmetics() bool {
	return true
}
