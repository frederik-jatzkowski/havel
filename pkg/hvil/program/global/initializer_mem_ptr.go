package global

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type InitializerMemPtr struct {
	tool.Node[InitializerMemPtr]
	names.NameResolution[struct {
		Decl *Decl
	}]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Name string `parser:"'mem' '.' 'ptr' '(' @Ident ')'"`
}

func (node *InitializerMemPtr) ResolveNames(ctx context.Context) error {
	decl, err := contexttool.FromCtx[instruction.VarDecl](ctx, node.Name)
	if err != nil {
		return node.Wrap(err)
	}

	global, ok := decl.(*Decl)
	if !ok {
		return node.Errorf("declaration should be a global variable but was %T", decl)
	}

	node.NameResolutionPass.Decl = global

	return nil
}

func (node *InitializerMemPtr) ResolveTypes(expected types.Type) error {
	if err := (&types.Ref{}).EqualsDetailed(expected); err != nil {
		return node.Wrap(err)
	}

	node.TypeCheckPass.Type = expected

	return nil
}

func (node *InitializerMemPtr) GenerateVirtualMachineAssembly(p *assembly.P) error {
	p.AddSLit(node.TypeCheckPass.Type.Bytes(), virtualmachine.NewFatPtr(0, uint32(node.NameResolutionPass.Decl.RelAddr())).ToUint64())

	return nil
}
