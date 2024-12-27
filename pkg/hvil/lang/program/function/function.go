package function

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Function struct {
	tool.Node[Function]
	names.NameResolution[struct {
		Blocks names.Scope[block.Block]
		Vars   names.Scope[memory.VarDecl]
	}]

	Name   string                    `parser:"'func':Keyword @Ident"`
	Params tool.List[memory.VarDecl] `parser:"'(' @@ ')'"`
	Result memory.VarDecl            `parser:"( '=>' '(' @@ ')' )?"`
	Locals tool.List[memory.VarDecl] `parser:"'{' ( 'declare':Keyword '(' @@ ')' ';' )?"`
	Blocks []block.Block             `parser:"@@+  '}'"`
}

func (f Function) Identifier() string {
	return f.Name
}

func (f *Function) ResolveNames() (errs []error) {
	f.NameResolutionPass.Vars = names.NewRootScope[memory.VarDecl]("variable")

	errs = f.NameResolutionPass.Vars.DefineAll(f.Params.Items)
	errs = append(errs, f.NameResolutionPass.Vars.DefineAll(f.Locals.Items)...)

	if f.Result != nil {
		err := f.NameResolutionPass.Vars.Define(&f.Result)
		if err != nil {
			errs = append(errs, err)
		}
	}

	f.NameResolutionPass.Blocks = names.NewRootScope[block.Block]("block")
	errs = append(errs, f.NameResolutionPass.Blocks.DefineAll(f.Blocks)...)

	for i := 0; i < len(f.Blocks); i++ {
		errs = append(errs, f.Blocks[i].ResolveNames(f.NameResolutionPass.Vars)...)
	}

	_, exists := f.NameResolutionPass.Blocks.Find("entry")
	if !exists {
		errs = append(errs, f.Errorf("no entry block defined"))
	}

	return errs
}

func (f *Function) ResolveTypes() (errs []error) {
	for i := 0; i < len(f.Blocks); i++ {
		errs = append(errs, f.Blocks[i].ResolveTypes()...)
	}

	return errs
}
