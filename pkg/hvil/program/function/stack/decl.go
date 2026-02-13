package stack

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Decl struct {
	tool.Node[Decl]
	statistics.Statistics[struct {
		PtrTaken   bool
		Reads      map[statistics.BlockID][]statistics.InstructionID
		Writes     map[statistics.BlockID][]statistics.InstructionID
		LiveRanges map[statistics.BlockID][]controlflow.LiveRange
	}]
	address.Resolution[struct {
		RelAddr int
	}]

	Name         string     `parser:"@Ident"`
	DeclaredType types.Type `parser:"':' @@"`
}

func (node *Decl) Identifier() string {
	return node.Name
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

func (node *Decl) CalculateStatistics(_ context.Context, entry controlflow.Node) {
	node.StatisticsPass.LiveRanges = controlflow.ComputeLiveRanges(entry, node.StatisticsPass.Reads, node.StatisticsPass.Writes)
}

func (node *Decl) RelAddr() int {
	return node.AddressResolutionPass.RelAddr
}

func (node *Decl) SetPtrTaken() {
	node.StatisticsPass.PtrTaken = true
}

func (node *Decl) AddBytecodeVirtualmachinePtrInstruction(p *assembly.P, target bytecode.R, dereferences []uint) error {
	_, offset, err := node.DeclaredType.Dereference(dereferences)
	if err != nil {
		return err
	}

	p.AddI1RLit(bytecode.OPStackPtr, target, uint16(node.RelAddr()+int(offset)), node.Position())

	return nil
}
