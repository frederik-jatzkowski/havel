package statements

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Content interface {
	tool.NodeLike
	names.Resolver
	typecheck.Resolver
}

type Statement struct {
	tool.Node[Statement]

	Content Content `parser:"@@ ';'"`
}

func (node *Statement) ResolveNames(ctx context.Context) error {
	return node.Content.ResolveNames(ctx)
}

func (node *Statement) ResolveTypes(ctx context.Context) error {
	return node.Content.ResolveTypes(ctx)
}
