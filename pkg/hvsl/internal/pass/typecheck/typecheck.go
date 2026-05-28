package typecheck

import (
	"context"
)

type TypeCheck[T any] struct {
	TypeCheckPass T `parser:"" json:",omitempty"`
}

type Resolver interface {
	ResolveTypes(ctx context.Context) error
}
