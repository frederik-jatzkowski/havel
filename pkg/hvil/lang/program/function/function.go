package function

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Function struct {
	tool.Node[Function]
	names.NameResolution[struct {
		Entry  *block.Block
		Blocks names.Scope[block.Block]
		Vars   names.Scope[memory.VarDecl]
	}]
	address.Resolution[struct {
		FrameSize  int
		ArgsSize   int
		ReturnSize int
		VarsSize   int
		RegsSize   int
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

	entry, exists := f.NameResolutionPass.Blocks.Find("entry")
	if !exists {
		errs = append(errs, f.Errorf("no entry block defined"))
	}

	f.NameResolutionPass.Entry = entry

	return errs
}

func (f *Function) ResolveTypes() (errs []error) {
	for i := 0; i < len(f.Blocks); i++ {
		errs = append(errs, f.Blocks[i].ResolveTypes()...)
	}

	return errs
}

func (f *Function) ResolveAddresses() (errs []error) {
	regsOffset := f.AddressResolutionPass.ArgsSize +
		f.AddressResolutionPass.ReturnSize +
		f.AddressResolutionPass.VarsSize

	for i := 0; i < len(f.Blocks); i++ {
		blockRegSize := 0
		for _, reg := range f.Blocks[i].NameResolutionPass.OrderedRegs {
			reg.AddressResolutionPass.RelAddr = regsOffset + blockRegSize
			blockRegSize += reg.Type().Bytes()
		}

		f.AddressResolutionPass.RegsSize = max(blockRegSize, f.AddressResolutionPass.RegsSize)
	}

	f.AddressResolutionPass.FrameSize += f.AddressResolutionPass.RegsSize

	return errs
}

func (f *Function) Execute(vm *runtime.VirtualMachine) (err error) {
	next := f.NameResolutionPass.Entry
	for next != nil {
		next, err = next.Execute(vm)
		if err != nil {
			return err
		}
	}

	return nil
}
