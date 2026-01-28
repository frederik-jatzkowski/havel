package stack

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
)

type Decl struct {
	tool.Node[Decl]
	statistics.Statistics[struct {
		Reads      map[statistics.BlockID][]statistics.InstructionID
		Writes     map[statistics.BlockID][]statistics.InstructionID
		LiveRanges map[statistics.BlockID][]controlflow.LiveRange
	}]
	address.Resolution[struct {
		RelAddr int
	}]
	registeralloc.RegisterAllocation[struct {
		BoundTo architecture.Register
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

func (node *Decl) CalculateStatistics(_ context.Context, entry controlflow.Node) {
	node.StatisticsPass.LiveRanges = controlflow.ComputeLiveRanges(entry, node.StatisticsPass.Reads, node.StatisticsPass.Writes)
}
