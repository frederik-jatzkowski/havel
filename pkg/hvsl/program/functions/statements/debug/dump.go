package debug

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Dump struct {
	tool.Node[Dump]

	Param statements.Expression `parser:"'dump' '(' @@ ')'"`
}

func (node *Dump) ResolveNames(ctx context.Context) error {
	return node.Param.ResolveNames(ctx)
}

func (node *Dump) ResolveTypes(ctx context.Context) error {
	return node.Param.ResolveTypes(ctx)
}
