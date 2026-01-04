package function

import (
	"context"

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

func (node *Function) ResolveNames(ctx context.Context) error {
	node.NameResolutionPass.Vars = names.NewRootScope[*stack.Decl](names.KindVariable)
	ctx = stack.WithScope(ctx, node.NameResolutionPass.Vars)

	for _, param := range node.Params.Items {
		if err := node.NameResolutionPass.Vars.Define(param); err != nil {
			return param.Wrap(err)
		}
	}

	if node.Result != nil {
		if err := node.NameResolutionPass.Vars.Define(node.Result); err != nil {
			return node.Result.Wrap(err)
		}
	}

	for _, local := range node.Locals.Items {
		if err := node.NameResolutionPass.Vars.Define(local); err != nil {
			return local.Wrap(err)
		}
	}

	node.NameResolutionPass.Blocks = names.NewRootScope[*block.Block](names.KindBlock)
	ctx = block.WithScope(ctx, node.NameResolutionPass.Blocks)

	for _, b := range node.Blocks {
		if err := node.NameResolutionPass.Blocks.Define(b); err != nil {
			return b.Wrap(err)
		}
	}

	for i := 0; i < len(node.Blocks); i++ {
		if err := node.Blocks[i].ResolveNames(ctx); err != nil {
			return err
		}
	}

	entry, err := node.NameResolutionPass.Blocks.Find(names.SpecialEntry)
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
