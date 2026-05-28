package program

import (
	"context"
	"errors"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/program"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type Program struct {
	tool.Node[Program]
	names.NameResolution[struct {
		Scope scope.Scope[scope.Object]
	}]

	Members []Member `parser:"@@*"`
}

func (node *Program) ResolveNames(ctx context.Context) error {
	node.NameResolutionPass.Scope, ctx = contexttool.WithScope[scope.Object](ctx, names.KindIdentifier)

	if err := errors.Join(
		contexttool.DefineInScope[scope.Object](ctx, types.BuiltinBool.Decl()),
		contexttool.DefineInScope[scope.Object](ctx, types.BuiltinU8.Decl()),
		contexttool.DefineInScope[scope.Object](ctx, types.BuiltinU16.Decl()),
		contexttool.DefineInScope[scope.Object](ctx, types.BuiltinU32.Decl()),
		contexttool.DefineInScope[scope.Object](ctx, types.BuiltinU64.Decl()),
	); err != nil {
		return node.Wrap(err)
	}

	for _, member := range node.Members {
		if object, ok := member.(scope.Object); ok {
			if err := contexttool.DefineInScope[scope.Object](ctx, object); err != nil {
				return err
			}
		}
	}

	for _, member := range node.Members {
		if err := member.ResolveNames(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (node *Program) ResolveTypes(ctx context.Context) error {
	for i := range node.Members {
		if err := node.Members[i].ResolveTypes(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (node *Program) GenerateHVILProgram() *program.Program {
	panic("implement me")
	return nil
}
