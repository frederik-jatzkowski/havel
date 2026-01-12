package contexttool

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type contextKeyScope[T any] struct{}

func WithScope[T names.ScopedObject](ctx context.Context, scope names.Scope[T]) context.Context {
	return context.WithValue(ctx, contextKeyScope[T]{}, scope)
}

func FromCtx[T names.ScopedObject](ctx context.Context, name string) (value T, err error) {
	scope, ok := ctx.Value(contextKeyScope[T]{}).(names.Scope[T])
	if !ok {
		return value, fmt.Errorf("no scope found in context for %s", reflect.TypeFor[T]())
	}

	return scope.Find(name)
}

type contextKeyCurrent[T any] struct{}

func WithCurrent[T any](ctx context.Context, current T) context.Context {
	return context.WithValue(ctx, contextKeyCurrent[T]{}, current)
}

func CurrentFromContext[T any](ctx context.Context) (T, error) {
	current, ok := ctx.Value(contextKeyCurrent[T]{}).(T)
	if !ok {
		return current, errors.New("no current object found in context")
	}

	return current, nil
}
