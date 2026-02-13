package call

import (
	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
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

type savedRegisters struct {
	Register architecture.Register
	RelAddr  int
	Bytes    int
}

func calculateSavedMemory(
	current *function.Function,
	block *block.Block,
	instructionID statistics.InstructionID,
	scope registeralloc.Scope,
) []savedRegisters {
	toSave := make([]savedRegisters, 0)

	for regWrite := range block.RegisterScope().All() {
		r := regWrite.Register()
		if scope.IsLiveAt(r, instructionID) {
			toSave = append(toSave, savedRegisters{
				Register: r,
				RelAddr:  regWrite.AddressResolutionPass.RelAddr,
				Bytes:    regWrite.Type().Bytes(),
			})
		}
	}

	return toSave
}

func generateVirtualMachineAssemblySaveCode(
	node tool.NodeLike,
	temp bytecode.R,
	toSave []savedRegisters,
	p *assembly.P,
) error {
	for _, saved := range toSave {
		p.AddI1RLit(bytecode.OPStackPtr, temp, uint16(saved.RelAddr), node.Position())

		op, err := bytecode.StoreForSize(saved.Bytes)
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI2R(op, temp, saved.Register.(bytecode.R), node.Position())
	}

	return nil
}

func generateVirtualMachineAssemblyParamsCode(
	node tool.NodeLike,
	temp1, temp2 bytecode.R,
	frameSize int,
	args []instruction.MemoryRead,
	p *assembly.P,
) error {
	offset := 16
	for _, arg := range args {
		if err := arg.GenerateVirtualMachineAssembly(p); err != nil {
			return node.Wrap(err)
		}

		size := arg.Type().Bytes()

		if _, ok := arg.Type().(*types.Composite); !ok {
			p.AddI1RLit(bytecode.OPStackPtr, temp1, uint16(frameSize+offset), node.Position())

			op, err := bytecode.StoreForSize(size)
			if err != nil {
				return node.Wrap(err)
			}

			p.AddI2R(op, temp1, arg.Register().(bytecode.R), node.Position())

			offset += size

			continue
		}

		sourceBase := arg.Register().(bytecode.R) // here is the initial base offset
		p.AddI1RLit(bytecode.OPStackPtr, temp1, uint16(frameSize+offset), node.Position())
		p.AddLit(temp2, 4, uint64(size), node.Position())

		p.AddI3R(bytecode.OPCopy, temp1, sourceBase, temp2, node.Position())

		offset += size
	}

	return nil
}

func generateVirtualMachineAssemblyRestoreCode(
	node tool.NodeLike,
	temp bytecode.R,
	toSave []savedRegisters,
	p *assembly.P,
) error {
	for _, saved := range toSave {
		p.AddI1RLit(bytecode.OPStackPtr, temp, uint16(saved.RelAddr), node.Position())

		op, err := bytecode.LoadForSize(saved.Bytes)
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI2R(op, saved.Register.(bytecode.R), temp, node.Position())
	}

	return nil
}
