package local

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
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
		Temp   architecture.Register
		Result architecture.Register
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

func (node *Call) AllocateRegisters(arch architecture.Architecture) ([]architecture.Register, error) {
	temp, ok := arch.GetScratchRegister()
	if !ok {
		return nil, node.Wrap(fmt.Errorf("failed to allocate register"))
	}

	node.RegisterAllocationPass.Temp = temp

	for _, arg := range node.Args.Items {
		regs, err := arg.AllocateRegisters(arch)
		if err != nil {
			return nil, node.Wrap(err)
		}

		arch.ReturnScratchRegisters(regs...)
	}

	return []architecture.Register{temp}, nil
}

func (node *Call) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *Call) CalculateLiveRanges(ctx context.Context) error {
	id, err := contexttool.CurrentFromContext[liveness.InstructionID](ctx)
	if err != nil {
		return node.Wrap(err)
	}

	node.LivenessPass.InstructionID = id

	for _, arg := range node.Args.Items {
		if err := arg.CalculateLiveRanges(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (node *Call) GenerateVirtualMachineAssembly(p *assembly.P) error {
	toSave := make([]*memory.RegWrite, 0)

	for regWrite := range node.NameResolutionPass.Block.RegisterScope().All() {
		if regWrite.WasLiveBefore(node.LivenessPass.InstructionID) && regWrite.WillBeLiveAfter(node.LivenessPass.InstructionID) {
			toSave = append(toSave, regWrite)
		}
	}

	for _, regWrite := range toSave {
		op, err := bytecode.StoreStackForSize(regWrite.Type().Bytes())
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI1RLit(op, regWrite.Register().(bytecode.R), uint16(regWrite.AddressResolutionPass.RelAddr), node.Position())
	}

	frameSize := node.NameResolutionPass.Current.AddressResolutionPass.FrameSize

	for i, param := range node.NameResolutionPass.Called.Params.Items {
		arg := node.Args.Items[i]

		op, err := bytecode.StoreStackForSize(param.Type().Bytes())
		if err != nil {
			return node.Wrap(err)
		}

		if err := arg.GenerateVirtualMachineAssembly(p); err != nil {
			return node.Wrap(err)
		}

		p.AddI1RLit(op, arg.Register().(bytecode.R), uint16(frameSize+param.AddressResolutionPass.RelAddr), node.Position())
	}

	temp := node.RegisterAllocationPass.Temp.(bytecode.R)

	// advance stack pointer
	p.AddI1RLit(bytecode.OPStoreStack64, bytecode.SP, uint16(frameSize+8), node.Position())
	p.AddLit(temp, 2, uint64(frameSize), node.Position())
	p.AddI3R(bytecode.OPAluAddU64, bytecode.SP, bytecode.SP, node.RegisterAllocationPass.Temp.(bytecode.R), node.Position())

	// prepare return address
	p.AddLit(temp, 1, 2, node.Position())
	p.AddI3R(bytecode.OPAluAddU64, temp, bytecode.PC, temp, node.Position())
	p.AddI1RLit(bytecode.OPStoreStack64, temp, 0, node.Position())

	p.AddJumpToLabel(node.NameResolutionPass.Called.NameResolutionPass.Entry.FullyQualifiedIdentifier(), node.Position())

	// restore stack pointer
	p.AddI1RLit(bytecode.OPLoadStack64, bytecode.SP, 8, node.Position())

	// restore registers
	for _, regWrite := range toSave {
		op, err := bytecode.LoadStackForSize(regWrite.Type().Bytes())
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI1RLit(op, regWrite.Register().(bytecode.R), uint16(regWrite.AddressResolutionPass.RelAddr), node.Position())
	}

	// load result
	void := types.Void{}
	if !void.Equals(node.TypeCheckPass.Signature.ReturnValue) {
		op, err := bytecode.LoadStackForSize(node.NameResolutionPass.Called.Result.Type().Bytes())
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI1RLit(op, node.RegisterAllocationPass.Result.(bytecode.R), uint16(frameSize+node.NameResolutionPass.Called.Result.AddressResolutionPass.RelAddr), node.Position())
	}

	return nil
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
