package contexttool

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type contextKeyScope[T any] struct{}

func WithScope[T scope.Object](ctx context.Context, kind fmt.Stringer) (scope.Scope[T], context.Context) {
	parent, ok := ctx.Value(contextKeyScope[T]{}).(scope.Scope[T])
	if !ok {
		root := scope.NewRoot[T](kind)
		return root, context.WithValue(ctx, contextKeyScope[T]{}, root)
	}

	child := parent.Child()

	return child, context.WithValue(ctx, contextKeyScope[T]{}, child)
}

func DefineInScope[T scope.Object](ctx context.Context, node T) error {
	scope, ok := ctx.Value(contextKeyScope[T]{}).(scope.Scope[T])
	if !ok {
		return errors.New("no scope found in context")
	}

	return scope.Define(node)
}

func FromCtx[T scope.Object](ctx context.Context, name string) (value T, err error) {
	scope, ok := ctx.Value(contextKeyScope[T]{}).(scope.Scope[T])
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
