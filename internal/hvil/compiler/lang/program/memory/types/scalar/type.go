package scalar

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
	"strconv"
)

type Type struct {
	tool.Node[Type]

	BitSize uint8 `parser:"@BitSize"`
}

func (t Type) String() string {
	return strconv.FormatUint(uint64(t.BitSize), 10)
}

func (t Type) Equals(other types.Type) bool {
	return t.String() == other.String()
}
