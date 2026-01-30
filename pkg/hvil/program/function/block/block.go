package block

import (
	"context"
	"fmt"
	"math"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Block struct {
	tool.Node[Block]
	names.NameResolution[struct {
		Regs     names.Scope[*instruction.RegWrite]
		Function *function.Function
	}]
	statistics.Statistics[struct {
		BlockID            statistics.BlockID
		FirstInstructionID statistics.InstructionID
		LastInstructionID  statistics.InstructionID
	}]
	address.Resolution[struct {
		SpillAddressMap map[*instruction.RegWrite]uint16
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
	node.NameResolutionPass.Regs = names.NewRootScope[*instruction.RegWrite](names.KindRegister)
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
	node.AddressResolutionPass.SpillAddressMap = make(map[*instruction.RegWrite]uint16)
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

func (node *Block) RegisterScope() names.Scope[*instruction.RegWrite] {
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

func (node *Block) CalculateStatistics(
	ctx context.Context,
	blockID statistics.BlockID,
	current statistics.InstructionID,
) (next statistics.InstructionID) {
	node.StatisticsPass.BlockID = blockID
	ctx = contexttool.WithCurrent(ctx, blockID)
	node.NameResolutionPass.Function.StatisticsPass.BlockCount++
	node.StatisticsPass.FirstInstructionID = current

	for i, _ := range node.Instructions {
		instr := &node.Instructions[i]
		ctx = contexttool.WithCurrent(ctx, current)
		node.NameResolutionPass.Function.StatisticsPass.InstructionCount++
		instr.CalculateStatistics(ctx)

		current++
	}

	ctx = contexttool.WithCurrent(ctx, current)
	node.Terminator.CalculateStatistics(ctx)
	node.StatisticsPass.LastInstructionID = current

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

	scope.SetInstructionID(node.StatisticsPass.LastInstructionID)

	if err := node.Terminator.AllocateRegisters(scope); err != nil {
		return err
	}

	scope.ReturnGeneralPurposeRegisters(registers...)

	return nil
}

func (node *Block) GenerateVirtualMachineAssembly(p *assembly.P) error {
	p.AddLabel(node.FullyQualifiedIdentifier(), node.Position())

	temp := node.NameResolutionPass.Function.RegisterAllocationPass.Temp.(bytecode.R)
	for _, param := range node.NameResolutionPass.Function.Params.Items {
		if param.RegisterAllocationPass.Volatile {
			p.AddI1RLit(bytecode.OPStackPtr, temp, uint16(param.AddressResolutionPass.RelAddr), node.Position())

			op, err := bytecode.StoreForSize(param.Type().Bytes())
			if err != nil {
				return node.Wrap(err)
			}

			p.AddI2R(op, temp, param.RegisterAllocationPass.BoundTo.(bytecode.R), node.Position())
		}
	}

	for i := range node.Instructions {
		if err := node.Instructions[i].GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}
	}

	return node.Terminator.GenerateVirtualMachineAssembly(p)
}

var _ controlflow.Node = (*Block)(nil)

func (node *Block) ID() statistics.BlockID {
	return node.StatisticsPass.BlockID
}

func (node *Block) Successors() []controlflow.Node {
	return node.Terminator.Successors()
}
