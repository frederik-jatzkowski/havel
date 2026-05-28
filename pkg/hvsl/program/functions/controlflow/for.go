package controlflow

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type For struct {
	tool.Node[For]

	Cond    statements.Expression `parser:"'for':Keyword @@ ';'"`
	After   statements.Content    `parser:"@@?"`
	Members []functions.Member    `parser:"'{' @@* '}'"`
}

func (node *For) ResolveNames(ctx context.Context) error {
	if err := node.Cond.ResolveNames(ctx); err != nil {
		return err
	}

	if node.After != nil {
		if err := node.After.ResolveNames(ctx); err != nil {
			return err
		}
	}

	for i := range node.Members {
		if err := node.Members[i].ResolveNames(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (node *For) ResolveTypes(ctx context.Context) error {
	if err := node.Cond.ResolveTypes(ctx); err != nil {
		return err
	}

	if node.After != nil {
		if err := node.After.ResolveTypes(ctx); err != nil {
			return err
		}
	}

	for i := range node.Members {
		if err := node.Members[i].ResolveTypes(ctx); err != nil {
			return err
		}
	}

	return nil
}
