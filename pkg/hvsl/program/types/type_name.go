package types

import (
	"context"
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type TypeName struct {
	tool.Node[TypeName]
	names.NameResolution[struct {
		Decl *Decl
	}]
	typecheck.TypeCheck[struct {
		ResolvedType Type
	}]

	Ptr  []string `parser:"@'*'*"`
	Name string   `parser:"@Ident"`
}

func (node *TypeName) ResolveNames(ctx context.Context) error {
	object, err := contexttool.FromCtx[scope.Object](ctx, node.Name)
	if err != nil {
		return node.Wrap(err)
	}

	t, ok := object.(*Decl)
	if !ok {
		return node.Wrap(fmt.Errorf("%s object is not a type declaration", node.Name))
	}

	node.NameResolutionPass.Decl = t

	return nil
}

func (node *TypeName) ResolveTypes(ctx context.Context) error {
	if err := node.NameResolutionPass.Decl.ResolveTypes(ctx); err != nil {
		return err
	}

	node.TypeCheckPass.ResolvedType = node.resolveType()

	return nil
}

func (node *TypeName) Type() Type {
	if node == nil {
		return &Void{}
	}

	if node.TypeCheckPass.ResolvedType == nil {
		node.TypeCheckPass.ResolvedType = node.resolveType()
	}

	return node.TypeCheckPass.ResolvedType
}

func (node *TypeName) resolveType() Type {
	if node == nil {
		return &Void{}
	}

	result := node.NameResolutionPass.Decl.Typ
	for range node.Ptr {
		result = &Ptr{Underlying: result}
	}

	return result
}
