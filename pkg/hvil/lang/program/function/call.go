package function

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Call struct {
	tool.Node[Call]
	names.NameResolution[struct {
		Decl *Function
	}]
	tool.NotImplemented[Call]

	Name string                 `parser:"'local' '.' @Ident"`
	Args tool.List[memory.Read] `parser:"'(' @@ ')'"`
}

func (c *Call) ResolveNames(vars names.Scope[memory.VarDecl], regs names.Scope[memory.RegWrite]) (errs []error) {
	return nil
}

func (c *Call) ResolveTypes(target types.Type) (errs []error) {
	return append(errs, c.Errorf("not implemented"))
}
