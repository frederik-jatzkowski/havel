package block

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type Block struct {
	tool.Node[Block]
	names.NameResolution[struct {
		Regs     names.Scope[*memory.RegWrite]
		Function *function.Function
	}]

	Name         string                    `parser:"'block':Keyword @Ident '{'"`
	Instructions []instruction.Instruction `parser:"@@*"`
	Terminator   Terminator                `parser:"'}' @@ ';'"`
}

func (node *Block) Identifier() string {
	return node.Name
}

func (node *Block) FullyQualifiedIdentifier() string {
	return node.NameResolutionPass.Function.Identifier() + "." + node.Identifier()
}

func (node *Block) ResolveNames(ctx context.Context) error {
	node.NameResolutionPass.Regs = names.NewRootScope[*memory.RegWrite](names.KindRegister)
	ctx = memory.WithRegisterScope(ctx, node.NameResolutionPass.Regs)

	for _, i := range node.Instructions {
		if err := i.ResolveNames(ctx); err != nil {
			return err
		}
	}

	fn, err := contexttool.CurrentFromContext[*function.Function](ctx)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Function = fn

	return node.Terminator.ResolveNames(ctx)
}

func (node *Block) RegisterScope() names.Scope[*memory.RegWrite] {
	return node.NameResolutionPass.Regs
}

func (node *Block) ResolveTypes() error {
	for i := 0; i < len(node.Instructions); i++ {
		if err := node.Instructions[i].ResolveTypes(); err != nil {
			return err
		}
	}

	return node.Terminator.ResolveTypes()
}

func (node *Block) AllocateRegisters(arch architecture.Architecture) error {
	registers := make([]architecture.Register, 0, len(node.Instructions))
	for i := range node.Instructions {
		regs, err := node.Instructions[i].AllocateRegisters(arch)
		if err != nil {
			return err
		}

		registers = append(registers, regs...)
	}

	for _, reg := range registers {
		arch.ReturnGeneralPurposeRegisters(reg)
	}

	return node.Terminator.AllocateRegisters(arch)
}

func (node *Block) GenerateVirtualMachineAssembly(p *assembly.P) error {
	p.AddLabel(node.FullyQualifiedIdentifier(), node.Position())

	for i := range node.Instructions {
		if err := node.Instructions[i].GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}
	}

	return node.Terminator.GenerateVirtualMachineAssembly(p)
}

func (node *Block) Execute(vm *runtime.VirtualMachine) (function.Block, error) {
	for _, i := range node.Instructions {
		err := i.Execute(vm)
		if err != nil {
			return nil, err
		}
	}

	return node.Terminator.Execute(vm)
}
