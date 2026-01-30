package types

import (
	"fmt"
	"strconv"

	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Scalar struct {
	tool.Node[Scalar]

	Size int `parser:"@Size? 'byte'"`
}

func (node *Scalar) String() string {
	return strconv.FormatUint(uint64(node.Size), 10) + " byte"
}

func (node *Scalar) MarshalText() ([]byte, error) {
	return []byte(node.String()), nil
}

func (node *Scalar) CanBeAssigned(other Type) bool {
	return node.Equals(other)
}

func (node *Scalar) CanBeAssignedDetailed(other Type) error {
	return node.EqualsDetailed(other)
}

func (node *Scalar) Equals(other Type) bool {
	return node.EqualsDetailed(other) == nil
}

func (node *Scalar) EqualsDetailed(other Type) error {
	otherScalar, ok := other.(*Scalar)
	if !ok {
		return fmt.Errorf("%s is not a scalar type", other)
	}

	if node.Bytes() != otherScalar.Bytes() {
		return fmt.Errorf("expected %d bytes, got %d", node.Bytes(), otherScalar.Bytes())
	}

	return nil
}

func (node *Scalar) Bytes() int {
	if node.Size == 0 {
		return 1
	}

	return node.Size
}

func (node *Scalar) CanDoArithmetics() bool {
	return true
}
