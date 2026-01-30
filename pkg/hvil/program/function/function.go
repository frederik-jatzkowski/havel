package function

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
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
		AddressTaken      int
	}]
	address.Resolution[struct {
		FrameSize int
		VarsSize  int
		RegsSize  int
		CallPlan  architecture.CallPlan
	}]
	registeralloc.RegisterAllocation[struct {
		Temp architecture.Register
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
	ctx = contexttool.WithScope(ctx, node.NameResolutionPass.Vars)

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

func (node *Function) Signature() *types.Function {
	signature := &types.Function{
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

func (node *Function) CalculateStatistics(ctx context.Context) {
	var current statistics.InstructionID
	var blockID statistics.BlockID
	for _, block := range node.Blocks {
		current = block.CalculateStatistics(ctx, blockID, current)
		blockID++
	}

	for _, param := range node.Params.Items {
		param.CalculateStatistics(ctx, node.NameResolutionPass.Entry)
	}

	if node.Result != nil {
		node.Result.CalculateStatistics(ctx, node.NameResolutionPass.Entry)
	}

	for _, local := range node.Locals.Items {
		local.CalculateStatistics(ctx, node.NameResolutionPass.Entry)
	}
}

func (node *Function) ResolveAddresses(arch architecture.Architecture) error {
	callPlan := arch.CalculateCallPlan(node.Signature())
	node.AddressResolutionPass.CallPlan = callPlan

	for i, paramPlan := range callPlan.Params {
		param := node.Params.Items[i]
		param.AddressResolutionPass.RelAddr = paramPlan.RelAddr
		param.RegisterAllocationPass.Volatile = node.Params.Items[i].StatisticsPass.PtrTaken
		param.RegisterAllocationPass.BoundTo = paramPlan.BoundTo
	}

	if node.Result != nil {
		node.Result.AddressResolutionPass.RelAddr = callPlan.Result.RelAddr
		node.Result.RegisterAllocationPass.Volatile = node.Result.StatisticsPass.PtrTaken
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
	temp, ok := scope.GetScratchRegister()
	if !ok {
		return node.Errorf("failed to allocate temp register")
	}

	node.RegisterAllocationPass.Temp = temp
	scope.ReturnScratchRegisters(temp)

	for i := 0; i < len(node.Blocks); i++ {
		if err := node.Blocks[i].AllocateRegisters(scope); err != nil {
			return err
		}
	}

	return nil
}

func (node *Function) GenerateVirtualMachineAssembly(p *assembly.P) error {
	temp := node.RegisterAllocationPass.Temp.(bytecode.R)
	for _, param := range node.Params.Items {
		if param.RegisterAllocationPass.Volatile {
			p.AddI1RLit(bytecode.OPStackPtr, temp, uint16(param.AddressResolutionPass.RelAddr), node.Position())

			op, err := bytecode.StoreForSize(param.Type().Bytes())
			if err != nil {
				return node.Wrap(err)
			}

			p.AddI2R(op, temp, param.RegisterAllocationPass.BoundTo.(bytecode.R), node.Position())
		}
	}

	for i := 0; i < len(node.Blocks); i++ {
		if err := node.Blocks[i].GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}
	}

	return nil
}
