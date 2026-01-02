package function

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Function struct {
	tool.Node[Function]
	names.NameResolution[struct {
		Entry  *block.Block
		Blocks names.Scope[*block.Block]
		Vars   names.Scope[*stack.Decl]
	}]
	address.Resolution[struct {
		FrameSize  int
		ArgsSize   int
		ReturnSize int
		VarsSize   int
		RegsSize   int
	}]

	Name   string                 `parser:"'func':Keyword @Ident"`
	Params tool.List[*stack.Decl] `parser:"'(' @@ ')'"`
	Result *stack.Decl            `parser:"( '=>' '(' @@ ')' )?"`
	Locals tool.List[*stack.Decl] `parser:"'{' ( 'declare':Keyword '(' @@ ')' ';' )?"`
	Blocks []*block.Block         `parser:"@@+  '}'"`
}

func (node *Function) Identifier() string {
	return node.Name
}

func (node *Function) ResolveNames() error {
	node.NameResolutionPass.Vars = names.NewRootScope[*stack.Decl]("variable")

	if err := node.NameResolutionPass.Vars.DefineAll(node.Params.Items); err != nil {
		return err
	}

	if node.Result != nil {
		if err := node.NameResolutionPass.Vars.Define(node.Result); err != nil {
			return node.Result.Wrap(err)
		}
	}

	if err := node.NameResolutionPass.Vars.DefineAll(node.Locals.Items); err != nil {
		return err
	}

	node.NameResolutionPass.Blocks = names.NewRootScope[*block.Block]("block")
	if err := node.NameResolutionPass.Blocks.DefineAll(node.Blocks); err != nil {
		return err
	}

	for i := 0; i < len(node.Blocks); i++ {
		if err := node.Blocks[i].ResolveNames(node.NameResolutionPass.Vars); err != nil {
			return err
		}
	}

	entry, err := node.NameResolutionPass.Blocks.Find("entry")
	if err != nil {
		return node.Errorf("no entry block defined")
	}

	node.NameResolutionPass.Entry = entry

	return err
}

func (node *Function) ResolveTypes() error {
	for i := 0; i < len(node.Blocks); i++ {
		if err := node.Blocks[i].ResolveTypes(); err != nil {
			return err
		}
	}

	return nil
}

func (node *Function) ResolveAddresses() error {
	offset := 0
	node.resolveLocalsAddresses(offset)
	offset += node.AddressResolutionPass.VarsSize
	node.resolveRegisterAddresses(offset)

	node.AddressResolutionPass.FrameSize = offset

	return nil
}

func (node *Function) resolveLocalsAddresses(offset int) {
	for _, decl := range node.Locals.Items {
		size := decl.Type().Bytes()
		decl.AddressResolutionPass.RelAddr = offset + node.AddressResolutionPass.VarsSize
		node.AddressResolutionPass.VarsSize += size
	}
}

func (node *Function) resolveRegisterAddresses(offset int) {
	for i := 0; i < len(node.Blocks); i++ {
		blockRegSize := 0
		for _, reg := range node.Blocks[i].NameResolutionPass.OrderedRegs {
			reg.AddressResolutionPass.RelAddr = offset + blockRegSize
			blockRegSize += reg.Type().Bytes()
		}

		node.AddressResolutionPass.RegsSize = max(blockRegSize, node.AddressResolutionPass.RegsSize)
	}
}

func (node *Function) Execute(vm *runtime.VirtualMachine) (err error) {
	next := node.NameResolutionPass.Entry
	for next != nil {
		next, err = next.Execute(vm)
		if err != nil {
			return err
		}
	}

	return nil
}
