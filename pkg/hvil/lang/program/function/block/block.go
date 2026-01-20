package block

import (
	"context"
	"fmt"
	"math"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type Block struct {
	tool.Node[Block]
	names.NameResolution[struct {
		Regs     names.Scope[*memory.RegWrite]
		Function *function.Function
	}]
	statistics.Statistics[struct {
		TerminatorInstructionID statistics.InstructionID
	}]
	address.Resolution[struct {
		SpillAddressMap map[*memory.RegWrite]uint16
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
	ctx = contexttool.WithScope(ctx, node.NameResolutionPass.Regs)
	ctx = contexttool.WithCurrent(ctx, node)

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

func (node *Block) ResolveAddresses(offset int) int {
	node.AddressResolutionPass.SpillAddressMap = make(map[*memory.RegWrite]uint16)
	blockRegSize := 0
	for reg := range node.RegisterScope().All() {
		reg.AddressResolutionPass.RelAddr = offset + blockRegSize
		if reg.AddressResolutionPass.RelAddr > math.MaxUint16 {
			panic(fmt.Sprintf("address out of range: %d", reg.AddressResolutionPass.RelAddr))
		}

		blockRegSize += reg.Type().Bytes()
	}

	return blockRegSize
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

func (node *Block) CalculateStatistics(ctx context.Context, current statistics.InstructionID) (next statistics.InstructionID) {
	node.NameResolutionPass.Function.StatisticsPass.BlockCount++

	for i, _ := range node.Instructions {
		instr := &node.Instructions[i]
		ctx = contexttool.WithCurrent(ctx, current)
		node.NameResolutionPass.Function.StatisticsPass.InstructionCount++
		instr.CalculateStatistics(ctx)

		current++
	}

	ctx = contexttool.WithCurrent(ctx, current)
	node.Terminator.CalculateStatistics(ctx)
	node.StatisticsPass.TerminatorInstructionID = current

	return current + 1
}

func (node *Block) AllocateRegisters(scope registeralloc.Scope) error {
	registers := make([]architecture.Register, 0, len(node.Instructions))
	for i := range node.Instructions {
		regs, err := node.Instructions[i].AllocateRegisters(scope)
		if err != nil {
			return err
		}

		registers = append(registers, regs...)
	}

	scope.SetInstructionID(node.StatisticsPass.TerminatorInstructionID)

	if err := node.Terminator.AllocateRegisters(scope); err != nil {
		return err
	}

	scope.ReturnGeneralPurposeRegisters(registers...)

	return nil
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
