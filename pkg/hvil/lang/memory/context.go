package memory

import (
	"context"
	"errors"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type contextRegisterKey struct{}

func WithRegisterScope(ctx context.Context, scope names.Scope[*RegWrite]) context.Context {
	return context.WithValue(ctx, contextRegisterKey{}, scope)
}

func DefineRegisterInCtx(ctx context.Context, node *RegWrite) error {
	scope, ok := ctx.Value(contextRegisterKey{}).(names.Scope[*RegWrite])
	if !ok {
		return errors.New("no register scope found in context")
	}

	return scope.Define(node)
}

func RegisterFromCtx(ctx context.Context, name string) (*RegWrite, error) {
	scope, ok := ctx.Value(contextRegisterKey{}).(names.Scope[*RegWrite])
	if !ok {
		return nil, errors.New("no register scope found in context")
	}

	return scope.Find(name)
}
