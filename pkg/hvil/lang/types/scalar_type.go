package types

import (
	"strconv"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type ScalarType struct {
	tool.Node[ScalarType]

	Size int `parser:"@Size? 'byte'"`
}

func (node *ScalarType) String() string {
	if node.Size <= 1 {
		return "byte"
	}

	return strconv.FormatUint(uint64(node.Size), 10) + " byte"
}

func (node *ScalarType) MarshalText() ([]byte, error) {
	return []byte(node.String()), nil
}

func (node *ScalarType) CanBeAssigned(other Type) bool {
	return node.Bytes() == other.Bytes()
}

func (node *ScalarType) Bytes() int {
	if node.Size == 0 {
		return 1
	}

	return node.Size
}
