package function

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type Function struct {
	tool.Node[Function]
	names.NameResolution[struct {
		Entry  Block
		Blocks names.Scope[Block]
		Vars   names.Scope[*stack.Decl]
	}]
	statistics.Statistics[struct {
		BlockCount        int
		InstructionCount  int
		AddressTakenCount int
		CalledCount       int
	}]
	address.Resolution[struct {
		FrameSize int
		VarsSize  int
		RegsSize  int
	}]

	Name   string                 `parser:"'func':Keyword @Ident"`
	Params tool.List[*stack.Decl] `parser:"'(' @@ ')'"`
	Result *stack.Decl            `parser:"( '->' '(' @@ ')' )?"`
	Locals tool.List[*stack.Decl] `parser:"'{' ( 'declare':Keyword '(' @@ ')' ';' )?"`
	Blocks []Block                `parser:"@@+  '}'"`
}

func (node *Function) Identifier() string {
	return node.Name
}

func (node *Function) ResolveNames(ctx context.Context) error {
	ctx = contexttool.WithCurrent(ctx, node)

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

	node.NameResolutionPass.Blocks = names.NewRootScope[Block](names.KindBlock)
	ctx = contexttool.WithScope(ctx, node.NameResolutionPass.Blocks)

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

func (node *Function) Signature() *types.FunctionType {
	signature := &types.FunctionType{
		Parameters: tool.List[types.Type]{
			Items: make([]types.Type, 0, len(node.Params.Items)),
		},
	}

	for _, item := range node.Params.Items {
		signature.Parameters.Items = append(signature.Parameters.Items, item.Type())
	}

	if node.Result == nil {
		signature.ReturnValue = &types.Void{}
	} else {
		signature.ReturnValue = node.Result.Type()
	}

	return signature
}

func (node *Function) CalculateStatistics() {
	for _, block := range node.Blocks {
		block.CalculateStatistics()
	}

	for _, param := range node.Params.Items {
		param.CalculateStatistics()
	}

	node.Result.CalculateStatistics()

	for _, local := range node.Locals.Items {
		local.CalculateStatistics()
	}
}

func (node *Function) ResolveAddresses(arch architecture.Architecture) error {
	callPlan := arch.CalculateCallPlan(node.Signature())
	for i, paramPlan := range callPlan.Params {
		param := node.Params.Items[i]
		param.AddressResolutionPass.RelAddr = paramPlan.RelAddr
		param.RegisterAllocationPass.BoundTo = paramPlan.BoundTo
	}

	if node.Result != nil {
		node.Result.AddressResolutionPass.RelAddr = callPlan.Result.RelAddr
		node.Result.RegisterAllocationPass.BoundTo = callPlan.Result.BoundTo
	}

	offset := callPlan.Offset
	node.resolveLocalsAddresses(offset)
	offset += node.AddressResolutionPass.VarsSize
	node.resolveRegisterAddresses(offset)
	offset += node.AddressResolutionPass.RegsSize

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
		node.AddressResolutionPass.RegsSize = max(
			node.AddressResolutionPass.RegsSize,
			node.Blocks[i].ResolveAddresses(offset),
		)
	}
}

func (node *Function) AllocateRegisters(allocator registeralloc.Allocator) error {
	scope := allocator.NewScope()
	for i := 0; i < len(node.Blocks); i++ {
		if err := node.Blocks[i].AllocateRegisters(scope); err != nil {
			return err
		}
	}

	return nil
}

func (node *Function) GenerateVirtualMachineAssembly(p *assembly.P) error {
	for i := 0; i < len(node.Blocks); i++ {
		if err := node.Blocks[i].GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}
	}

	return nil
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
