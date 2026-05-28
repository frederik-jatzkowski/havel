package alu

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type UnOp struct {
	tool.Node[UnOp]
	typecheck.TypeCheck[struct {
		ResolvedType types.Type
	}]

	Symbol UnOpSymbol            `parser:"@@"`
	Arg    statements.Expression `parser:"'(' @@ ')'"`
}

func (node *UnOp) ResolveNames(ctx context.Context) error {
	if err := node.Arg.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *UnOp) ResolveTypes(ctx context.Context) error {
	if err := node.Arg.ResolveTypes(ctx); err != nil {
		return err
	}

	switch node.Symbol {
	case UnOpSymbolNot:
		if !node.Arg.Type().Equals(types.BuiltinBool) {
			return node.Errorf("argument must be of type bool, got %s", node.Arg.Type())
		}

		node.TypeCheckPass.ResolvedType = types.BuiltinBool
	}

	return nil
}

func (node *UnOp) Type() types.Type {
	return node.TypeCheckPass.ResolvedType
}
