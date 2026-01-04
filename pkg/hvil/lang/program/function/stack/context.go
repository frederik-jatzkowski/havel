package stack

import (
	"context"
	"errors"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type contextKey struct{}

func WithScope(ctx context.Context, scope names.Scope[*Decl]) context.Context {
	return context.WithValue(ctx, contextKey{}, scope)
}

func FromCtx(ctx context.Context, name string) (*Decl, error) {
	scope, ok := ctx.Value(contextKey{}).(names.Scope[*Decl])
	if !ok {
		return nil, errors.New("no local declaration scope found in context")
	}

	return scope.Find(name)
}
