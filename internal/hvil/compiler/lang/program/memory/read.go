package memory

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types"
)

type Read interface {
	Type() types.Type
	Position() lexer.Position
}
