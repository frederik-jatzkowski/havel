package debug

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Operation statements.Content

type Call struct {
	tool.Node[Call]

	Operation Operation `parser:"'debug':Keyword '.' @@"`
}

func (node *Call) ResolveNames(ctx context.Context) error {
	return node.Operation.ResolveNames(ctx)
}

func (node *Call) ResolveTypes(ctx context.Context) error {
	return node.Operation.ResolveTypes(ctx)
}
