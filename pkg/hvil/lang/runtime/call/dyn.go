package call

import (
	"context"
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Dyn struct {
	tool.Node[Dyn]
	names.NameResolution[struct {
		Current *function.Function
		Block   *block.Block
	}]
	typecheck.TypeCheck[struct {
		Signature *types.FunctionType
	}]
	statistics.Statistics[struct {
		InstructionID statistics.InstructionID
	}]
	registeralloc.RegisterAllocation[struct {
		Scope    registeralloc.Scope
		Temp     architecture.Register
		Result   architecture.Register
		CallPlan architecture.CallPlan
	}]

	Target instruction.MemoryRead            `parser:"'dyn' '.' @@"`
	Args   tool.List[instruction.MemoryRead] `parser:"'(' @@ ')'"`
}

func (node *Dyn) ResolveNames(ctx context.Context) error {
	if err := node.Target.ResolveNames(ctx); err != nil {
		return err
	}

	for _, item := range node.Args.Items {
		if err := item.ResolveNames(ctx); err != nil {
			return err
		}
	}

	var err error

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

func (node *Dyn) ResolveTypes(target types.Type) error {
	node.TypeCheckPass.Signature = calculateSignature(node.Args.Items, target)

	if err := node.TypeCheckPass.Signature.CanBeAssignedDetailed(node.Target.Type()); err != nil {
		return node.Wrap(err)
	}

	return nil
}

func (node *Dyn) CalculateStatistics(ctx context.Context) {
	instructionID, err := contexttool.CurrentFromContext[statistics.InstructionID](ctx)
	if err != nil {
		panic(err)
	}

	node.StatisticsPass.InstructionID = instructionID

	for _, arg := range node.Args.Items {
		arg.CalculateStatistics(ctx)
	}
}

func (node *Dyn) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	node.RegisterAllocationPass.Scope = scope
	node.RegisterAllocationPass.CallPlan = scope.Architecture().CalculateCallPlan(node.TypeCheckPass.Signature)

	temp, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Wrap(fmt.Errorf("failed to allocate register"))
	}

	node.RegisterAllocationPass.Temp = temp

	regs, err := node.Target.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	for _, arg := range node.Args.Items {
		regs, err := arg.AllocateRegisters(scope)
		if err != nil {
			return nil, node.Wrap(err)
		}

		scope.ReturnScratchRegisters(regs...)
	}

	scope.ReturnScratchRegisters(regs...)

	return []architecture.Register{temp}, nil
}

func (node *Dyn) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *Dyn) GenerateVirtualMachineAssembly(p *assembly.P) error {
	toSave := calculateSavedMemory(
		node.NameResolutionPass.Current,
		node.NameResolutionPass.Block,
		node.StatisticsPass.InstructionID,
		node.RegisterAllocationPass.Scope,
	)

	temp := node.RegisterAllocationPass.Temp.(bytecode.R)

	if err := generateVirtualMachineAssemblySaveCode(node, temp, toSave, p); err != nil {
		return err
	}

	frameSize := node.NameResolutionPass.Current.AddressResolutionPass.FrameSize
	callPlan := node.RegisterAllocationPass.CallPlan

	if err := generateVirtualMachineAssemblyParamsCode(node, temp, frameSize, callPlan, node.Args.Items, p); err != nil {
		return err
	}

	p.AddCall(node.Target.Register().(bytecode.R), uint32(frameSize), node.Position())

	// restore stack pointer
	p.AddI1RLit(bytecode.OPStackPtr, temp, 8, node.Position())
	p.AddI2R(bytecode.OPLoad64, bytecode.SP, temp, node.Position())

	result := node.RegisterAllocationPass.Result

	if result != nil {
		if err := generateVirtualMachineAssemblyResultCode(node, temp, frameSize, callPlan, result.(bytecode.R), p); err != nil {
			return err
		}
	}

	if err := generateVirtualMachineAssemblyRestoreCode(node, temp, toSave, p); err != nil {
		return err
	}

	return nil
}
