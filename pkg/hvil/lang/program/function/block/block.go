package block

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Block struct {
	tool.Node[Block]
	names.NameResolution[struct {
		Regs        names.Scope[memory.RegWrite]
		OrderedRegs []*memory.RegWrite
	}]

	Name         string                    `parser:"'block':Keyword @Ident '{'"`
	Instructions []instruction.Instruction `parser:"@@*"`
	Terminator   Terminator                `parser:"'}' '=>' @@ ';'"`
}

func (b Block) Identifier() string {
	return b.Name
}

func (b *Block) ResolveNames(vars names.Scope[memory.VarDecl]) (errs []error) {
	b.NameResolutionPass.Regs = names.NewRootScope[memory.RegWrite]("register")

	for _, i := range b.Instructions {
		errs = append(errs, i.ResolveNames(vars, b.NameResolutionPass.Regs)...)
		if reg, ok := i.Result.(*memory.RegWrite); ok {
			b.NameResolutionPass.OrderedRegs = append(b.NameResolutionPass.OrderedRegs, reg)
		}
	}

	return errs
}

func (b *Block) ResolveTypes() (errs []error) {
	for i := 0; i < len(b.Instructions); i++ {
		errs = append(errs, b.Instructions[i].ResolveTypes()...)
	}

	return errs
}

func (b *Block) Execute(vm *runtime.VirtualMachine) (*Block, error) {
	for _, i := range b.Instructions {
		err := i.Execute(vm)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
