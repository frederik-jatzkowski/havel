package global

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type Block struct {
	tool.Node[Block]
	names.NameResolution[struct {
		Declarations scope.Scope[instruction.VarDecl]
	}]
	address.Resolution[struct {
		RelAddr int
	}]

	Declarations tool.List[Decl] `parser:"'declare':Keyword '(' @@ ')' ';'"`
}

func (node *Block) ResolveNames(ctx context.Context) error {
	for i := range node.Declarations.Items {
		if err := node.Declarations.Items[i].ResolveNames(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (node *Block) ResolveTypes() error {
	for i := range node.Declarations.Items {
		if err := node.Declarations.Items[i].ResolveTypes(); err != nil {
			return err
		}
	}

	return nil
}

func (node *Block) ResolveAddresses(arch architecture.Architecture) error {
	offset := 0
	for i := range node.Declarations.Items {
		decl := &node.Declarations.Items[i]
		decl.AddressResolutionPass.RelAddr = offset
		offset += decl.Type().Bytes()
	}

	return nil
}

func (node *Block) GenerateVirtualMachineAssembly(p *assembly.P) error {
	for _, decl := range node.Declarations.Items {
		if err := decl.GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}
	}

	return nil
}
