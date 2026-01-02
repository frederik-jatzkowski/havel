package block

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Block struct {
	tool.Node[Block]
	names.NameResolution[struct {
		Regs        names.Scope[*memory.RegWrite]
		OrderedRegs []*memory.RegWrite
	}]

	Name         string                    `parser:"'block':Keyword @Ident '{'"`
	Instructions []instruction.Instruction `parser:"@@*"`
	Terminator   Terminator                `parser:"'}' '=>' @@ ';'"`
}

func (node *Block) Identifier() string {
	return node.Name
}

func (node *Block) ResolveNames(vars names.Scope[*stack.Decl]) error {
	node.NameResolutionPass.Regs = names.NewRootScope[*memory.RegWrite]("register")

	for _, i := range node.Instructions {
		if err := i.ResolveNames(vars, node.NameResolutionPass.Regs); err != nil {
			return err
		}

		if reg, ok := i.Result.(*memory.RegWrite); ok {
			node.NameResolutionPass.OrderedRegs = append(node.NameResolutionPass.OrderedRegs, reg)
		}
	}

	return nil
}

func (node *Block) ResolveTypes() error {
	for i := 0; i < len(node.Instructions); i++ {
		if err := node.Instructions[i].ResolveTypes(); err != nil {
			return err
		}
	}

	return nil
}

func (node *Block) Execute(vm *runtime.VirtualMachine) (*Block, error) {
	for _, i := range node.Instructions {
		err := i.Execute(vm)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
