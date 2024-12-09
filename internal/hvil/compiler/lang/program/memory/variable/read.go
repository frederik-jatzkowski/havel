package variable

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/pass"
)

type Read struct {
	tool.Node[Read]
	pass.NameResolution[struct {
		Decl Decl
	}]
	tool.NotImplemented[Read]

	Ident string `parser:"@Ident"`
}

func (read *Read) Type() types.Type {
	return read.NameResolutionPass.Decl.Type()
}

func (read *Read) Position() lexer.Position {
	return read.Pos
}
