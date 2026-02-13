package call

import (
	"context"
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Local struct {
	tool.Node[Local]
	names.NameResolution[struct {
		Current *function.Function
		Block   *block.Block
		Called  *function.Function
	}]
	typecheck.TypeCheck[struct {
		Signature *types.Function
	}]
	statistics.Statistics[struct {
		InstructionID statistics.InstructionID
	}]
	registeralloc.RegisterAllocation[struct {
		Scope        registeralloc.Scope
		Temp1, Temp2 architecture.Register
	}]

	Name string                            `parser:"'local' '.' @Ident"`
	Args tool.List[instruction.MemoryRead] `parser:"'(' @@ ')'"`
}

func (node *Local) ResolveNames(ctx context.Context) error {
	decl, err := contexttool.FromCtx[*function.Function](ctx, node.Name)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Called = decl

	for _, item := range node.Args.Items {
		if err := item.ResolveNames(ctx); err != nil {
			return err
		}
	}

	node.NameResolutionPass.Current, err = contexttool.CurrentFromContext[*function.Function](ctx)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Block, err = contexttool.CurrentFromContext[*block.Block](ctx)
	if err != nil {
		return node.Wrap(err)
	}

	return nil
}

func (node *Local) ResolveTypes(target types.Type) error {
	if !target.Equals(&types.Void{}) {
		return node.Errorf("cannot assign void to %s", target)
	}

	node.TypeCheckPass.Signature = calculateSignature(node.Args.Items)

	if err := node.TypeCheckPass.Signature.EqualsDetailed(node.NameResolutionPass.Called.Signature()); err != nil {
		return node.Wrap(err)
	}

	return nil
}

func (node *Local) CalculateStatistics(ctx context.Context) {
	instructionID, err := contexttool.CurrentFromContext[statistics.InstructionID](ctx)
	if err != nil {
		panic(err)
	}

	node.StatisticsPass.InstructionID = instructionID

	for _, arg := range node.Args.Items {
		arg.CalculateStatistics(ctx)
	}
}

func (node *Local) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	node.RegisterAllocationPass.Scope = scope

	temp1, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Wrap(fmt.Errorf("failed to allocate register"))
	}

	node.RegisterAllocationPass.Temp1 = temp1

	temp2, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Wrap(fmt.Errorf("failed to allocate register"))
	}

	node.RegisterAllocationPass.Temp2 = temp2

	for _, arg := range node.Args.Items {
		regs, err := arg.AllocateRegisters(scope)
		if err != nil {
			return nil, node.Wrap(err)
		}

		scope.ReturnScratchRegisters(regs...)
	}

	return []architecture.Register{temp1, temp2}, nil
}

func (node *Local) SetResultRegister(r architecture.Register) {
	panic(fmt.Sprintf("target register assigned to %T, which returns void", node))
}

func (node *Local) GenerateVirtualMachineAssembly(p *assembly.P) error {
	toSave := calculateSavedMemory(
		node.NameResolutionPass.Current,
		node.NameResolutionPass.Block,
		node.StatisticsPass.InstructionID,
		node.RegisterAllocationPass.Scope,
	)

	temp1 := node.RegisterAllocationPass.Temp1.(bytecode.R)
	temp2 := node.RegisterAllocationPass.Temp2.(bytecode.R)

	p.AddI0R(bytecode.OPDebugStackPush, node.Position())

	if err := generateVirtualMachineAssemblySaveCode(node, temp1, toSave, p); err != nil {
		return err
	}

	frameSize := node.NameResolutionPass.Current.AddressResolutionPass.FrameSize

	if err := generateVirtualMachineAssemblyParamsCode(node, temp1, temp2, frameSize, node.Args.Items, p); err != nil {
		return err
	}

	p.AddLoadLabel(temp1, node.NameResolutionPass.Called.NameResolutionPass.Entry.FullyQualifiedIdentifier(), node.Position())
	p.AddCall(temp1, uint32(frameSize), node.Position())

	// restore stack pointer
	p.AddI1RLit(bytecode.OPStackPtr, temp1, 8, node.Position())
	p.AddI2R(bytecode.OPLoad64, bytecode.SP, temp1, node.Position())

	if err := generateVirtualMachineAssemblyRestoreCode(node, temp1, toSave, p); err != nil {
		return err
	}

	p.AddI0R(bytecode.OPDebugStackPop, node.Position())

	return nil
}
