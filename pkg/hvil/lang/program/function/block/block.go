package block

import (
	"context"
	"errors"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type Block struct {
	tool.Node[Block]
	names.NameResolution[struct {
		Regs        names.Scope[*memory.RegWrite]
		OrderedRegs []*memory.RegWrite
	}]

	Name         string                    `parser:"'block':Keyword @Ident '{'"`
	Instructions []instruction.Instruction `parser:"@@*"`
	Terminator   Terminator                `parser:"'}' @@ ';'"`
}

func (node *Block) Identifier() string {
	return node.Name
}

func (node *Block) ResolveNames(ctx context.Context) error {
	node.NameResolutionPass.Regs = names.NewRootScope[*memory.RegWrite](names.KindRegister)
	ctx = memory.WithRegisterScope(ctx, node.NameResolutionPass.Regs)

	for _, i := range node.Instructions {
		if err := i.ResolveNames(ctx); err != nil {
			return err
		}

		if reg, ok := i.Result().(*memory.RegWrite); ok {
			node.NameResolutionPass.OrderedRegs = append(node.NameResolutionPass.OrderedRegs, reg)
		}
	}

	return node.Terminator.ResolveNames(ctx)
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
	for i := range node.Instructions {
		instr := &node.Instructions[i]

		regWrite, ok := instr.Result().(*memory.RegWrite)
		if !ok {
			continue
		}

		r, ok := arch.GetGeneralPurposeRegister()
		if !ok {
			return errors.New("no general purpose registers remaining")
		}

		regWrite.RegisterAllocationPass.Register = r
		instr.Operation.SetResultRegister(r)
	}

	return nil
}

func (node *Block) GenerateVirtualMachineAssembly(p *assembly.P, isMain bool) error {
	for i := range node.Instructions {
		if err := node.Instructions[i].GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}
	}

	return node.Terminator.GenerateVirtualMachineAssembly(p, isMain)
}

func (node *Block) Execute(vm *runtime.VirtualMachine) (*Block, error) {
	for _, i := range node.Instructions {
		err := i.Execute(vm)
		if err != nil {
			return nil, err
		}
	}

	return node.Terminator.Execute(vm)
}
