package types

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Decl struct {
	tool.Node[Decl]

	Name string `parser:"'type':Keyword @Ident"`
	Typ  Type   `parser:"@@"`
}

func (node *Decl) Identifier() string {
	return node.Name
}

func (node *Decl) ResolveNames(ctx context.Context) error {
	ctx = contexttool.WithCurrent(ctx, node)

	return node.Typ.ResolveNames(ctx)
}

func (node *Decl) ResolveTypes(ctx context.Context) error {
	return nil
}

func (node *Decl) Type() Type {
	return node.Typ
}
