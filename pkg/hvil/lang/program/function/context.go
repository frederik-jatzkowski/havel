package function

import (
	"context"
	"errors"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type contextKey struct{}

func WithScope(ctx context.Context, scope names.Scope[*Function]) context.Context {
	return context.WithValue(ctx, contextKey{}, scope)
}

func FromCtx(ctx context.Context, name string) (*Function, error) {
	scope, ok := ctx.Value(contextKey{}).(names.Scope[*Function])
	if !ok {
		return nil, errors.New("no functions scope found in context")
	}

	return scope.Find(name)
}
