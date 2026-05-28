package alu

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Operation interface {
	statements.Content
	types.TypedObject
}

type Call struct {
	tool.Node[Call]

	Operation Operation `parser:"'alu':Keyword '.' @@"`
}

func (node *Call) ResolveNames(ctx context.Context) error {
	return node.Operation.ResolveNames(ctx)
}

func (node *Call) ResolveTypes(ctx context.Context) error {
	return node.Operation.ResolveTypes(ctx)
}

func (node *Call) Type() types.Type {
	return node.Operation.Type()
}
