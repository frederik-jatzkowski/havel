package alu

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type BinOp struct {
	tool.Node[BinOp]
	typecheck.TypeCheck[struct {
		Input  types.Type
		Output types.Type
	}]

	Symbol BinOpSymbol           `parser:"@@"`
	Left   statements.Expression `parser:"'(' @@ ','"`
	Right  statements.Expression `parser:"@@ ')'"`
}

func (node *BinOp) ResolveNames(ctx context.Context) error {
	if err := node.Left.ResolveNames(ctx); err != nil {
		return err
	}

	if err := node.Right.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *BinOp) ResolveTypes(ctx context.Context) error {
	if err := node.Left.ResolveTypes(ctx); err != nil {
		return err
	}

	if err := node.Right.ResolveTypes(ctx); err != nil {
		return err
	}

	if !node.Left.Type().Equals(node.Right.Type()) {
		return node.Errorf("operands have different types: %s vs. %s", node.Left.Type(), node.Right.Type())
	}

	if err := node.validateAllowedInput(); err != nil {
		return err
	}

	node.TypeCheckPass.Input = node.Left.Type()
	node.TypeCheckPass.Output = node.determineOutputType()

	return nil
}

func (node *BinOp) validateAllowedInput() error {
	builtin, ok := node.Left.Type().(types.Builtin)
	if !ok {
		return node.Errorf("expected builtin type, got %s", node.Left.Type())
	}

	switch node.Symbol {
	case
		BinOpSymbolAdd,
		BinOpSymbolSub,
		BinOpSymbolMul,
		BinOpSymbolDiv,
		BinOpSymbolMod,
		BinOpSymbolLeq,
		BinOpSymbolLt,
		BinOpSymbolGeq,
		BinOpSymbolGt:
		switch builtin {
		case
			types.BuiltinU8,
			types.BuiltinU16,
			types.BuiltinU32,
			types.BuiltinU64:
			return nil
		default:
			return node.Errorf("invalid type for %s: %s", node.Symbol, builtin)
		}
	}

	return nil
}

func (node *BinOp) determineOutputType() types.Type {
	switch node.Symbol {
	case BinOpSymbolAdd, BinOpSymbolSub, BinOpSymbolMul, BinOpSymbolDiv, BinOpSymbolMod:
		return node.Left.Type()
	case BinOpSymbolLeq, BinOpSymbolLt, BinOpSymbolGeq, BinOpSymbolGt:
		return types.BuiltinBool
	case BinOpSymbolEq:
		return types.BuiltinBool
	default:
		return nil
	}
}

func (node *BinOp) Type() types.Type {
	return node.TypeCheckPass.Output
}
