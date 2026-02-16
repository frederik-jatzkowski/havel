package global

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
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
	Initializer  Initializer `parser:"( '=' @@ )?"`
}

func (node *Decl) Identifier() string {
	return node.Name
}

func (node *Decl) ResolveNames(ctx context.Context) error {
	if node.Initializer != nil {
		if err := node.Initializer.ResolveNames(ctx); err != nil {
			return err
		}
	}

	return contexttool.DefineInScope[instruction.VarDecl](ctx, node)
}

func (node *Decl) ResolveTypes() error {
	if node.Initializer != nil {
		return node.Initializer.ResolveTypes(node.DeclaredType)
	}

	return nil
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

func (node *Decl) RelAddr() int {
	return node.AddressResolutionPass.RelAddr
}

func (node *Decl) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if node.Initializer != nil {
		return node.Initializer.GenerateVirtualMachineAssembly(p)
	}

	// fill with zeros
	size := node.DeclaredType.Bytes()
	for size > 0 {
		switch {
		case size >= 8:
			p.AddSLit(8, 0)
			size -= 8
		case size >= 4:
			p.AddSLit(4, 0)
			size -= 4
		case size >= 2:
			p.AddSLit(2, 0)
			size -= 2
		default:
			p.AddSLit(1, 0)
			size -= 1
		}
	}

	return nil
}

func (node *Decl) AddBytecodeVirtualmachinePtrInstruction(p *assembly.P, target bytecode.R, dereferences []uint) error {
	_, offset, err := node.DeclaredType.Dereference(dereferences)
	if err != nil {
		return err
	}

	p.AddI1RLit(bytecode.OPStaticPtr, target, uint16(node.RelAddr()+int(offset)), node.Position())

	return nil
}
