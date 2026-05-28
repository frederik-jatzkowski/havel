package controlflow

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type If struct {
	tool.Node[If]

	Cond    statements.Expression `parser:"'if':Keyword @@"`
	Members []functions.Member    `parser:"'{' (@@)* '}'"`
	Else    *Else                 `parser:"@@?"`
}

func (node *If) ResolveNames(ctx context.Context) error {
	if err := node.Cond.ResolveNames(ctx); err != nil {
		return err
	}

	for i := range node.Members {
		if err := node.Members[i].ResolveNames(ctx); err != nil {
			return err
		}
	}

	if node.Else != nil {
		return node.Else.ResolveNames(ctx)
	}

	return nil
}

func (node *If) ResolveTypes(ctx context.Context) error {
	if err := node.Cond.ResolveTypes(ctx); err != nil {
		return err
	}

	if !node.Cond.Type().Equals(types.BuiltinBool) {
		return node.Errorf("condition must be of type bool, got %s", node.Cond.Type())
	}

	for i := range node.Members {
		if err := node.Members[i].ResolveTypes(ctx); err != nil {
			return err
		}
	}

	if node.Else != nil {
		return node.Else.ResolveTypes(ctx)
	}

	return nil
}
