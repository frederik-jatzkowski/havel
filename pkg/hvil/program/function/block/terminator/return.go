package terminator

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Return struct {
	tool.Node[Return]
	names.NameResolution[struct {
		IsMain   bool
		Function *function.Function
	}]
	registeralloc.RegisterAllocation[struct {
		Temp architecture.Register
	}]

	Token string `parser:"@'return':Keyword" json:"-"`
}

var _ block.Terminator = (*Return)(nil)

func (node *Return) ResolveNames(ctx context.Context) error {
	fn, err := contexttool.CurrentFromContext[*function.Function](ctx)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.IsMain = fn.Identifier() == names.SpecialMain
	node.NameResolutionPass.Function = fn

	return nil
}

func (node *Return) ResolveTypes() error {
	return nil
}

func (node *Return) CalculateStatistics(ctx context.Context) {
	result := node.NameResolutionPass.Function.Result
	if result != nil {
		blockID, err := contexttool.CurrentFromContext[statistics.BlockID](ctx)
		if err != nil {
			panic(err)
		}
		instructionID, err := contexttool.CurrentFromContext[statistics.InstructionID](ctx)
		if err != nil {
			panic(err)
		}

		if result.StatisticsPass.Reads == nil {
			result.StatisticsPass.Reads = make(map[statistics.BlockID][]statistics.InstructionID)
		}

		result.StatisticsPass.Reads[blockID] = append(
			result.StatisticsPass.Reads[blockID],
			instructionID,
		)
	}
}

func (node *Return) AllocateRegisters(scope registeralloc.Scope) error {
	r, ok := scope.GetScratchRegister()
	if !ok {
		return node.Errorf("failed to obtain exit code register")
	}

	scope.ReturnScratchRegisters(r)
	node.RegisterAllocationPass.Temp = r

	return nil
}

func (node *Return) GenerateVirtualMachineAssembly(p *assembly.P) error {
	temp := node.RegisterAllocationPass.Temp.(bytecode.R)
	if node.NameResolutionPass.IsMain {
		// exit code 0
		p.AddLit(temp, 1, 0, node.Position())
		p.AddI1R(bytecode.OPExit, temp, node.Position())
	} else {
		result := node.NameResolutionPass.Function.Result
		if result != nil && result.RegisterAllocationPass.Volatile {
			// store final result in register if it is kept in memory

			p.AddI1RLit(bytecode.OPStackPtr, temp, uint16(result.AddressResolutionPass.RelAddr), node.Position())

			op, err := bytecode.LoadForSize(result.Type().Bytes())
			if err != nil {
				return node.Wrap(err)
			}

			p.AddI2R(op, result.RegisterAllocationPass.BoundTo.(bytecode.R), temp, node.Position())
		}

		p.AddI1RLit(bytecode.OPStackPtr, temp, 0, node.Position())
		p.AddI2R(bytecode.OPLoad64, bytecode.PC, temp, node.Position())
	}

	return nil
}

func (node *Return) Successors() []controlflow.Node {
	return nil
}
