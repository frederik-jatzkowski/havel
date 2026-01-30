package call

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

func calculateSignature(args []instruction.MemoryRead, target types.Type) *types.FunctionType {
	signature := &types.FunctionType{
		Parameters: tool.List[types.Type]{
			Items: make([]types.Type, 0, len(args)),
		},
	}

	for _, item := range args {
		signature.Parameters.Items = append(signature.Parameters.Items, item.Type())
	}

	signature.ReturnValue = target

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
		r := param.RegisterAllocationPass.BoundTo
		if r == nil || param.RegisterAllocationPass.Volatile {
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
			RelAddr: param.AddressResolutionPass.RelAddr,
			Bytes:   param.Type().Bytes(),
		})
	}

	result := current.Result
	if result != nil {
		if r := result.RegisterAllocationPass.BoundTo; r != nil && !result.RegisterAllocationPass.Volatile && controlflow.MustBeSavedAt(
			result.StatisticsPass.LiveRanges[blockID],
			instructionID,
		) {
			toSave = append(toSave, architecture.MemoryAllocation{
				BoundTo: r,
				RelAddr: result.AddressResolutionPass.RelAddr,
				Bytes:   result.Type().Bytes(),
			})
		}
	}

	for _, local := range current.Locals.Items {
		r := local.RegisterAllocationPass.BoundTo
		if r == nil || local.RegisterAllocationPass.Volatile {
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
			RelAddr: local.AddressResolutionPass.RelAddr,
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

func generateVirtualMachineAssemblyResultCode(
	node tool.NodeLike,
	temp bytecode.R,
	frameSize int,
	callPlan architecture.CallPlan,
	result bytecode.R,
	p *assembly.P,
) error {
	if callPlan.Result.Bytes > 0 {
		plan := callPlan.Result
		if plan.BoundTo != nil {
			if plan.BoundTo != result {
				p.AddI2R(bytecode.OPAluMove, result, plan.BoundTo.(bytecode.R), node.Position())
			}

			return nil
		}

		p.AddI1RLit(bytecode.OPStackPtr, temp, uint16(frameSize+callPlan.Result.RelAddr), node.Position())

		op, err := bytecode.LoadForSize(callPlan.Result.Bytes)
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI2R(op, result, temp, node.Position())
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
