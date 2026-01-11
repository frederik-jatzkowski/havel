package local

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type Call struct {
	tool.Node[Call]
	names.NameResolution[struct {
		Current *function.Function
		Called  *function.Function
	}]
	typecheck.TypeCheck[struct {
		Signature *types.FunctionType
	}]
	tool.NotImplemented[Call]

	Name string                 `parser:"'local' '.' @Ident"`
	Args tool.List[memory.Read] `parser:"'(' @@ ')'"`
}

func (node *Call) ResolveNames(ctx context.Context) error {
	decl, err := function.FromCtx(ctx, node.Name)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Called = decl

	for _, item := range node.Args.Items {
		if err := item.ResolveNames(ctx); err != nil {
			return err
		}
	}

	node.NameResolutionPass.Current, err = function.CurrentFromContext(ctx)
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
	//TODO implement me
	panic("implement me")
}

func (node *Call) SetResultRegister(r architecture.Register) {
	//TODO implement me
	panic("implement me")
}

func (node *Call) GenerateVirtualMachineAssembly(p *assembly.P) error {
	//TODO implement me
	panic("implement me")
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
