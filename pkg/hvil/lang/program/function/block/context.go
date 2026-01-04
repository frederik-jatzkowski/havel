package block

import (
	"context"
	"errors"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type contextKey struct{}

func WithScope(ctx context.Context, scope names.Scope[*Block]) context.Context {
	return context.WithValue(ctx, contextKey{}, scope)
}

func FromCtx(ctx context.Context, name string) (*Block, error) {
	scope, ok := ctx.Value(contextKey{}).(names.Scope[*Block])
	if !ok {
		return nil, errors.New("no block scope found in context")
	}

	return scope.Find(name)
}
