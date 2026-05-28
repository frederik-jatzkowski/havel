package mem

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Alloc struct {
	tool.Node[Alloc]

	Typ types.TypeName `parser:"'alloc' '(' @@ ')'"`
}

func (node *Alloc) ResolveNames(ctx context.Context) error {
	return node.Typ.ResolveNames(ctx)
}

func (node *Alloc) ResolveTypes(ctx context.Context) error {
	return node.Typ.ResolveTypes(ctx)
}

func (node *Alloc) Type() types.Type {
	return &types.Ptr{Underlying: node.Typ.Type()}
}
