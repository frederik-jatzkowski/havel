package local

import (
	"context"
	"fmt"
	"slices"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc/liveness"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Call struct {
	tool.Node[Call]
	names.NameResolution[struct {
		Current *function.Function
		Block   *block.Block
		Called  *function.Function
	}]
	typecheck.TypeCheck[struct {
		Signature *types.FunctionType
	}]
	registeralloc.RegisterAllocation[struct {
		Scope    registeralloc.Scope
		Temp     architecture.Register
		Result   architecture.Register
		CallPlan architecture.CallPlan
	}]
	liveness.Liveness[struct {
		InstructionID liveness.InstructionID
	}]

	Name string                            `parser:"'local' '.' @Ident"`
	Args tool.List[instruction.MemoryRead] `parser:"'(' @@ ')'"`
}

func (node *Call) ResolveNames(ctx context.Context) error {
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

func (node *Call) ResolveTypes(target types.Type) error {
	node.calculateSignature(target)

	if err := node.TypeCheckPass.Signature.CanBeAssignedDetailed(node.NameResolutionPass.Called.Signature()); err != nil {
		return node.Wrap(err)
	}

	return nil
}

func (node *Call) CalculateStatistics() {
	for _, arg := range node.Args.Items {
		arg.CalculateStatistics()
	}
}

func (node *Call) calculateSignature(target types.Type) {
	node.TypeCheckPass.Signature = &types.FunctionType{
		Parameters: tool.List[types.Type]{
			Items: make([]types.Type, 0, len(node.Args.Items)),
		},
	}

	for _, item := range node.Args.Items {
		node.TypeCheckPass.Signature.Parameters.Items = append(node.TypeCheckPass.Signature.Parameters.Items, item.Type())
	}

	node.TypeCheckPass.Signature.ReturnValue = target
}

func (node *Call) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	node.LivenessPass.InstructionID = scope.GetInstructionID()
	node.RegisterAllocationPass.Scope = scope
	node.RegisterAllocationPass.CallPlan = scope.Architecture().CalculateCallPlan(node.TypeCheckPass.Signature)

	temp, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Wrap(fmt.Errorf("failed to allocate register"))
	}

	node.RegisterAllocationPass.Temp = temp

	for _, arg := range node.Args.Items {
		regs, err := arg.AllocateRegisters(scope)
		if err != nil {
			return nil, node.Wrap(err)
		}

		scope.ReturnScratchRegisters(regs...)
	}

	return []architecture.Register{temp}, nil
}

func (node *Call) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *Call) GenerateVirtualMachineAssembly(p *assembly.P) error {
	toSave := node.calculateSavedMemory()

	if err := node.generateVirtualMachineAssemblySaveCode(p, toSave); err != nil {
		return err
	}

	if err := node.generateVirtualMachineAssemblyParamsCode(p); err != nil {
		return err
	}

	node.generateVirtualMachineAssemblyCallCode(p)

	if err := node.generateVirtualMachineAssemblyResultCode(p); err != nil {
		return err
	}

	if err := node.generateVirtualMachineAssemblyRestoreCode(p, toSave); err != nil {
		return err
	}

	return nil
}

func (node *Call) generateVirtualMachineAssemblyCallCode(p *assembly.P) {
	temp := node.RegisterAllocationPass.Temp.(bytecode.R)
	frameSize := node.NameResolutionPass.Current.AddressResolutionPass.FrameSize

	p.AddLoadLabel(temp, node.NameResolutionPass.Called.NameResolutionPass.Entry.FullyQualifiedIdentifier(), node.Position())
	p.AddCall(temp, uint32(frameSize), node.Position())
	p.AddI1RLit(bytecode.OPLoadStack64, bytecode.SP, 8, node.Position()) // restore stack pointer
}

func (node *Call) generateVirtualMachineAssemblySaveCode(p *assembly.P, toSave []architecture.MemoryAllocation) error {
	for _, saved := range toSave {
		op, err := bytecode.StoreStackForSize(saved.Bytes)
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI1RLit(op, saved.BoundTo.(bytecode.R), uint16(saved.RelAddr), node.Position())
	}

	return nil
}

func (node *Call) generateVirtualMachineAssemblyParamsCode(p *assembly.P) error {
	frameSize := node.NameResolutionPass.Current.AddressResolutionPass.FrameSize
	for i, param := range node.NameResolutionPass.Called.Params.Items {
		arg := node.Args.Items[i]

		if err := arg.GenerateVirtualMachineAssembly(p); err != nil {
			return node.Wrap(err)
		}

		plan := node.RegisterAllocationPass.CallPlan.Params[i]
		if plan.BoundTo != nil {
			if plan.BoundTo != arg.Register() {
				p.AddI2R(bytecode.OPAluMove, plan.BoundTo.(bytecode.R), arg.Register().(bytecode.R), arg.Position())
			}
		} else {
			op, err := bytecode.StoreStackForSize(param.Type().Bytes())
			if err != nil {
				return node.Wrap(err)
			}

			p.AddI1RLit(op, arg.Register().(bytecode.R), uint16(frameSize+param.AddressResolutionPass.RelAddr), node.Position())
		}
	}

	return nil
}

func (node *Call) generateVirtualMachineAssemblyRestoreCode(p *assembly.P, toSave []architecture.MemoryAllocation) error {
	for _, saved := range toSave {
		op, err := bytecode.LoadStackForSize(saved.Bytes)
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI1RLit(op, saved.BoundTo.(bytecode.R), uint16(saved.RelAddr), node.Position())
	}

	return nil
}

func (node *Call) generateVirtualMachineAssemblyResultCode(p *assembly.P) error {
	frameSize := node.NameResolutionPass.Current.AddressResolutionPass.FrameSize
	void := types.Void{}
	if !void.Equals(node.TypeCheckPass.Signature.ReturnValue) {
		plan := node.RegisterAllocationPass.CallPlan.Result
		if plan.BoundTo != nil {
			if plan.BoundTo != node.RegisterAllocationPass.Result {
				p.AddI2R(bytecode.OPAluMove, node.RegisterAllocationPass.Result.(bytecode.R), plan.BoundTo.(bytecode.R), node.Position())
			}

			return nil
		}

		op, err := bytecode.LoadStackForSize(node.NameResolutionPass.Called.Result.Type().Bytes())
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI1RLit(
			op,
			node.RegisterAllocationPass.Result.(bytecode.R),
			uint16(frameSize+node.NameResolutionPass.Called.Result.AddressResolutionPass.RelAddr),
			node.Position(),
		)
	}

	return nil
}

func (node *Call) calculateSavedMemory() []architecture.MemoryAllocation {
	toSave := make([]architecture.MemoryAllocation, 0)

	for _, param := range node.NameResolutionPass.Current.Params.Items {
		r := param.RegisterAllocationPass.BoundTo
		if r == nil {
			continue
		}

		toSave = append(toSave, architecture.MemoryAllocation{
			BoundTo: r,
			RelAddr: param.AddressResolutionPass.RelAddr,
			Bytes:   param.Type().Bytes(),
		})
	}

	result := node.NameResolutionPass.Current.Result
	if result != nil {
		if r := result.RegisterAllocationPass.BoundTo; r != nil {
			toSave = append(toSave, architecture.MemoryAllocation{
				BoundTo: r,
				RelAddr: result.AddressResolutionPass.RelAddr,
				Bytes:   result.Type().Bytes(),
			})
		}
	}

	for regWrite := range node.NameResolutionPass.Block.RegisterScope().All() {
		r := regWrite.Register()
		if node.RegisterAllocationPass.Scope.IsLiveAt(r, node.LivenessPass.InstructionID) {
			toSave = append(toSave, architecture.MemoryAllocation{
				BoundTo: r,
				RelAddr: regWrite.AddressResolutionPass.RelAddr,
				Bytes:   regWrite.Type().Bytes(),
			})
		}
	}

	return slices.DeleteFunc(toSave, func(allocation architecture.MemoryAllocation) bool {
		return allocation.BoundTo == node.RegisterAllocationPass.Result
	})
}

func (node *Call) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	newStackPointer := vm.StackPointer + node.NameResolutionPass.Current.AddressResolutionPass.FrameSize
	newStackSize := newStackPointer + node.NameResolutionPass.Called.AddressResolutionPass.FrameSize
	if len(vm.Stack) < newStackSize {
		return node.Errorf("stack overflow")
	}

	prevStackPtr := vm.StackPointer

	sourceAddresses := make([]unsafe.Pointer, len(node.Args.Items))
	for i, arg := range node.Args.Items {
		sourceAddresses[i] = arg.Addr(vm)
	}

	vm.StackPointer = newStackPointer

	for i, param := range node.NameResolutionPass.Called.Params.Items {
		length := param.Type().Bytes()
		addr := param.Addr(vm)
		copy(unsafe.Slice((*byte)(addr), length), unsafe.Slice((*byte)(sourceAddresses[i]), length))
	}

	if err := node.NameResolutionPass.Called.Execute(vm); err != nil {
		return err
	}

	if node.NameResolutionPass.Called.Result != nil && result != nil {
		length := node.NameResolutionPass.Called.Result.Type().Bytes()
		sourceAddr := node.NameResolutionPass.Called.Result.Addr(vm)
		copy(unsafe.Slice((*byte)(result), length), unsafe.Slice((*byte)(sourceAddr), length))
	}

	vm.StackPointer = prevStackPtr

	return nil
}
