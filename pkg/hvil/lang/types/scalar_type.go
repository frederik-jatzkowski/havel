package types

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"strconv"
)

type ScalarType struct {
	tool.Node[ScalarType]

	Size int `parser:"@Size? 'byte'"`
}

func (t ScalarType) String() string {
	if t.Size <= 1 {
		return "byte"
	}

	return strconv.FormatUint(uint64(t.Size), 10) + " byte"
}

func (t ScalarType) CanBeAssigned(other Type) bool {
	return t.BitSize() == other.BitSize()
}

func (t ScalarType) BitSize() int {
	if t.Size == 0 {
		return 8
	}

	return t.Size * 8
}
