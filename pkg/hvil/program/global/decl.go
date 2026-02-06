package global

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Decl struct {
	tool.Node[Decl]
	address.Resolution[struct {
		RelAddr int
	}]
	statistics.Statistics[struct {
		PtrTaken bool
		Reads    map[statistics.BlockID][]statistics.InstructionID
		Writes   map[statistics.BlockID][]statistics.InstructionID
	}]

	Name         string      `parser:"@Ident"`
	DeclaredType types.Type  `parser:"':' @@"`
	Initializer  Initializer `parser:"'=' @@"`
}

func (node *Decl) Identifier() string {
	return node.Name
}

func (node *Decl) ResolveNames(ctx context.Context) error {
	if err := node.Initializer.ResolveNames(ctx); err != nil {
		return err
	}

	return contexttool.DefineInScope[instruction.VarDecl](ctx, node)
}

func (node *Decl) ResolveTypes() error {
	return node.Initializer.ResolveTypes(node.DeclaredType)
}

func (node *Decl) Type() types.Type {
	return node.DeclaredType
}

func (node *Decl) AddReadToStatistic(blockID statistics.BlockID, instructionID statistics.InstructionID) {
	if node.StatisticsPass.Reads == nil {
		node.StatisticsPass.Reads = make(map[statistics.BlockID][]statistics.InstructionID)
	}

	node.StatisticsPass.Reads[blockID] = append(node.StatisticsPass.Reads[blockID], instructionID)
}

func (node *Decl) AddWriteToStatistic(blockID statistics.BlockID, instructionID statistics.InstructionID) {
	if node.StatisticsPass.Writes == nil {
		node.StatisticsPass.Writes = make(map[statistics.BlockID][]statistics.InstructionID)
	}

	node.StatisticsPass.Writes[blockID] = append(node.StatisticsPass.Writes[blockID], instructionID)
}

func (node *Decl) SetPtrTaken() {
	node.StatisticsPass.PtrTaken = true
}

func (node *Decl) BoundTo() architecture.Register {
	return nil
}

func (node *Decl) Volatile() bool {
	return true
}

func (node *Decl) RelAddr() int {
	return node.AddressResolutionPass.RelAddr
}

func (node *Decl) GenerateVirtualMachineAssembly(p *assembly.P) error {
	return node.Initializer.GenerateVirtualMachineAssembly(p)
}

func (node *Decl) AddBytecodeVirtualmachinePtrInstruction(p *assembly.P, target bytecode.R) {
	p.AddI1RLit(bytecode.OPStaticPtr, target, uint16(node.RelAddr()), node.Position())
}
