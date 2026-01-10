package memory

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type RegRead struct {
	tool.Node[RegRead]
	names.NameResolution[struct {
		Decl *RegWrite
	}]

	Ident string `parser:"'$' @Ident"`
}

func (node *RegRead) Identifier() string {
	return node.NameResolutionPass.Decl.Identifier()
}

func (node *RegRead) ResolveNames(ctx context.Context) error {
	decl, err := RegisterFromCtx(ctx, node.Ident)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Decl = decl

	return nil
}

func (node *RegRead) GenerateVirtualMachineAssembly(p *assembly.P) error {
	return nil
}

func (node *RegRead) Register() architecture.Register {
	return node.NameResolutionPass.Decl.Register()
}

func (node *RegRead) Type() types.Type {
	return node.NameResolutionPass.Decl.RegType
}

func (node *RegRead) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return node.NameResolutionPass.Decl.Addr(vm)
}
