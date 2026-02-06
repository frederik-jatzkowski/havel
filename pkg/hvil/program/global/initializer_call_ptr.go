package global

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type InitializerCallPtr struct {
	tool.Node[InitializerCallPtr]

	Name string `parser:"'call' '.' 'ptr' '(' @Ident ')'"`
}

func (node *InitializerCallPtr) ResolveNames(ctx context.Context) error {
	// TODO: implement me
	panic("implement me")
}

func (node *InitializerCallPtr) ResolveTypes(expected types.Type) error {
	//TODO implement me
	panic("implement me")
}

func (node *InitializerCallPtr) GenerateVirtualMachineAssembly(p *assembly.P) error {
	// TODO: implement me
	panic("implement me")
}
