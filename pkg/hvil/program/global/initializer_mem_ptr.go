package global

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type InitializerMemPtr struct {
	tool.Node[InitializerMemPtr]

	Name string `parser:"'mem' '.' 'ptr' '(' @Ident ')'"`
}

func (node *InitializerMemPtr) ResolveNames(ctx context.Context) error {
	// TODO: implement me
	panic("implement me")
}

func (node *InitializerMemPtr) ResolveTypes(expected types.Type) error {
	//TODO implement me
	panic("implement me")
}

func (node *InitializerMemPtr) GenerateVirtualMachineAssembly(p *assembly.P) error {
	// TODO: implement me
	panic("implement me")
}
