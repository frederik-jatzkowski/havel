package mem

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type ElementPtr struct {
	tool.Node[ElementPtr]
	typecheck.TypeCheck[struct {
		Type   types.Type
		Offset uint
		IsVoid bool
	}]
	registeralloc.RegisterAllocation[struct {
		Temp1,
		Temp2,
		Result architecture.Register
	}]

	Ptr          instruction.MemoryRead `parser:"'element_ptr' '(' @@ ','"`
	DeclaredType types.Type             `parser:"@@ ','"`
	ConstIndex   uint                   `parser:"( @Number |"`
	Index        instruction.MemoryRead `parser:"@@ ) ( ','"`
	Dereferences []uint                 `parser:"'self' ('[' @Number ']')* )? ')'"`
}

func (node *ElementPtr) ResolveNames(ctx context.Context) error {
	if err := node.Ptr.ResolveNames(ctx); err != nil {
		return err
	}

	if node.Index != nil {
		if err := node.Index.ResolveNames(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (node *ElementPtr) ResolveTypes(target types.Type) error {
	if _, ok := node.Ptr.Type().(*types.Ref); !ok {
		return node.Ptr.Errorf("%s is not a ref type", node.Ptr.Type())
	}

	switch target.(type) {
	case *types.Ref:
	case *types.Void:
		node.TypeCheckPass.IsVoid = true
	default:
		return node.Errorf("cannot assign ref to %s", target)
	}

	node.TypeCheckPass.Type = target

	_, offset, err := node.DeclaredType.Dereference(node.Dereferences)
	if err != nil {
		return err
	}

	node.TypeCheckPass.Offset = offset

	if node.Index == nil {
		node.TypeCheckPass.Offset += node.ConstIndex

		return nil
	}

	if err := (&types.Scalar{Size: 8}).EqualsDetailed(node.Index.Type()); err != nil {
		return node.Index.Wrap(err)
	}

	return nil
}

func (node *ElementPtr) CalculateStatistics(ctx context.Context) {
	node.Ptr.CalculateStatistics(ctx)

	if node.Index != nil {
		node.Index.CalculateStatistics(ctx)
	}
}

func (node *ElementPtr) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	ptrRegs, err := node.Ptr.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	var indexRegs []architecture.Register

	if node.Index != nil {
		indexRegs, err = node.Index.AllocateRegisters(scope)
		if err != nil {
			return nil, err
		}
	}

	temp1, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Errorf("cannot allocate temp registers")
	}

	temp2, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Errorf("cannot allocate temp registers")
	}

	node.RegisterAllocationPass.Temp1 = temp1
	node.RegisterAllocationPass.Temp2 = temp2

	return append(append([]architecture.Register{temp1, temp2}, ptrRegs...), indexRegs...), nil
}

func (node *ElementPtr) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *ElementPtr) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if !node.TypeCheckPass.IsVoid {
		if err := node.Ptr.GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}

		temp1 := node.RegisterAllocationPass.Temp1.(bytecode.R)
		temp2 := node.RegisterAllocationPass.Temp2.(bytecode.R)

		if node.Index == nil {
			p.AddLit(temp1, 8, uint64(node.TypeCheckPass.Offset), node.Position())
			p.AddI3R(bytecode.OPAluAddU64, node.RegisterAllocationPass.Result.(bytecode.R), node.Ptr.Register().(bytecode.R), temp1, node.Position())

			return nil
		}

		if err := node.Index.GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}

		p.AddLit(temp1, 8, uint64(node.DeclaredType.Bytes()), node.Position())
		p.AddI3R(bytecode.OPAluMulU64, temp1, temp1, node.Index.Register().(bytecode.R), node.Position())

		if node.TypeCheckPass.Offset > 0 {
			p.AddLit(temp2, 8, uint64(node.TypeCheckPass.Offset), node.Position())
			p.AddI3R(bytecode.OPAluAddU64, temp1, temp1, temp2, node.Position())
		}

		p.AddI3R(bytecode.OPAluAddU64, node.RegisterAllocationPass.Result.(bytecode.R), node.Ptr.Register().(bytecode.R), temp1, node.Position())
	}

	return nil
}
