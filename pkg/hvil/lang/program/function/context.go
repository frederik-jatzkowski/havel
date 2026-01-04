package function

import (
	"context"
	"errors"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type contextScopeKey struct{}

func WithScope(ctx context.Context, scope names.Scope[*Function]) context.Context {
	return context.WithValue(ctx, contextScopeKey{}, scope)
}

func FromCtx(ctx context.Context, name string) (*Function, error) {
	scope, ok := ctx.Value(contextScopeKey{}).(names.Scope[*Function])
	if !ok {
		return nil, errors.New("no functions scope found in context")
	}

	return scope.Find(name)
}

type contextCurrentKey struct{}

func WithCurrent(ctx context.Context, current *Function) context.Context {
	return context.WithValue(ctx, contextCurrentKey{}, current)
}

func CurrentFromContext(ctx context.Context) (*Function, error) {
	current, ok := ctx.Value(contextCurrentKey{}).(*Function)
	if !ok {
		return nil, errors.New("no functions scope found in context")
	}

	return current, nil
}
