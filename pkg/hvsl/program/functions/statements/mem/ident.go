package mem

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type Decl interface {
	types.TypedObject
}

type Ident struct {
	tool.Node[Ident]
	names.NameResolution[struct {
		Decl Decl
	}]

	Name string `parser:"@Ident"`
}

func (node *Ident) ResolveNames(ctx context.Context) error {
	obj, err := contexttool.FromCtx[scope.Object](ctx, node.Name)
	if err != nil {
		return node.Wrap(err)
	}

	decl, ok := obj.(Decl)
	if !ok {
		return node.Errorf("referenced object is not a memory declaration, got '%s'", obj)
	}

	node.NameResolutionPass.Decl = decl

	return nil
}

func (node *Ident) ResolveTypes(ctx context.Context) error {
	return nil
}

func (node *Ident) Type() types.Type {
	return node.NameResolutionPass.Decl.Type()
}
