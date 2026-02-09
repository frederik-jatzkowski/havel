package call

import (
	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

func calculateSignature(args []instruction.MemoryRead) *types.Function {
	signature := &types.Function{
		Parameters: tool.List[types.Type]{
			Items: make([]types.Type, 0, len(args)),
		},
	}

	for _, item := range args {
		signature.Parameters.Items = append(signature.Parameters.Items, item.Type())
	}

	return signature
}

func calculateSavedMemory(
	current *function.Function,
	block *block.Block,
	instructionID statistics.InstructionID,
	scope registeralloc.Scope,
) []architecture.MemoryAllocation {
	blockID := block.StatisticsPass.BlockID

	toSave := make([]architecture.MemoryAllocation, 0)

	for _, param := range current.Params.Items {
		r := param.BoundTo()
		if r == nil || param.Volatile() {
			continue
		}

		if !controlflow.MustBeSavedAt(
			param.StatisticsPass.LiveRanges[blockID],
			instructionID,
		) {
			continue
		}

		toSave = append(toSave, architecture.MemoryAllocation{
			BoundTo: r,
			RelAddr: param.RelAddr(),
			Bytes:   param.Type().Bytes(),
		})
	}

	for _, local := range current.Locals.Items {
		r := local.BoundTo()
		if r == nil || local.Volatile() {
			continue
		}

		if !controlflow.MustBeSavedAt(
			local.StatisticsPass.LiveRanges[blockID],
			instructionID,
		) {
			continue
		}

		toSave = append(toSave, architecture.MemoryAllocation{
			BoundTo: r,
			RelAddr: local.RelAddr(),
			Bytes:   local.Type().Bytes(),
		})
	}

	for regWrite := range block.RegisterScope().All() {
		r := regWrite.Register()
		if scope.IsLiveAt(r, instructionID) {
			toSave = append(toSave, architecture.MemoryAllocation{
				BoundTo: r,
				RelAddr: regWrite.AddressResolutionPass.RelAddr,
				Bytes:   regWrite.Type().Bytes(),
			})
		}
	}

	return toSave
}

func generateVirtualMachineAssemblySaveCode(
	node tool.NodeLike,
	temp bytecode.R,
	toSave []architecture.MemoryAllocation,
	p *assembly.P,
) error {
	for _, saved := range toSave {
		p.AddI1RLit(bytecode.OPStackPtr, temp, uint16(saved.RelAddr), node.Position())

		op, err := bytecode.StoreForSize(saved.Bytes)
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI2R(op, temp, saved.BoundTo.(bytecode.R), node.Position())
	}

	return nil
}

func generateVirtualMachineAssemblyParamsCode(
	node tool.NodeLike,
	temp bytecode.R,
	frameSize int,
	callPlan architecture.CallPlan,
	args []instruction.MemoryRead,
	p *assembly.P,
) error {
	for i, param := range callPlan.Params {
		arg := args[i]

		if err := arg.GenerateVirtualMachineAssembly(p); err != nil {
			return node.Wrap(err)
		}

		if param.BoundTo != nil {
			if param.BoundTo != arg.Register() {
				p.AddI2R(bytecode.OPAluMove, param.BoundTo.(bytecode.R), arg.Register().(bytecode.R), arg.Position())
			}
		} else {
			p.AddI1RLit(bytecode.OPStackPtr, temp, uint16(frameSize+param.RelAddr), node.Position())

			op, err := bytecode.StoreForSize(param.Bytes)
			if err != nil {
				return node.Wrap(err)
			}

			p.AddI2R(op, temp, arg.Register().(bytecode.R), node.Position())
		}
	}

	return nil
}

func generateVirtualMachineAssemblyRestoreCode(
	node tool.NodeLike,
	temp bytecode.R,
	toSave []architecture.MemoryAllocation,
	p *assembly.P,
) error {
	for _, saved := range toSave {
		p.AddI1RLit(bytecode.OPStackPtr, temp, uint16(saved.RelAddr), node.Position())

		op, err := bytecode.LoadForSize(saved.Bytes)
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI2R(op, saved.BoundTo.(bytecode.R), temp, node.Position())
	}

	return nil
}
