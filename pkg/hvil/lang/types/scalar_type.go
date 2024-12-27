package types

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"strconv"
)

type ScalarType struct {
	tool.Node[ScalarType]

	Size int `parser:"@Size"`
}

func (t ScalarType) String() string {
	return strconv.FormatUint(uint64(t.Size), 10)
}

func (t ScalarType) CanBeAssigned(other Type) bool {
	return t.BitSize() == other.BitSize()
}

func (t ScalarType) BitSize() int {
	return t.Size
}
