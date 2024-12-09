package function

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/pass"
)

type Call struct {
	tool.Node[Call]
	pass.NameResolution[struct {
		Decl *Function
	}]
	tool.NotImplemented[Call]

	Name string                 `parser:"'local' '.' @Ident"`
	Args tool.List[memory.Read] `parser:"'(' @@ ')'"`
}
