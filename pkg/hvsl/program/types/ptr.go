package types

import (
	"context"
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Ptr struct {
	tool.Node[Struct]

	Underlying Type
}

func (node *Ptr) ResolveNames(ctx context.Context) error {
	return node.Underlying.ResolveNames(ctx)
}

func (node *Ptr) Equals(other Type) bool {
	ptr, ok := other.(*Ptr)
	if !ok {
		return false
	}

	return node.Underlying.Equals(ptr.Underlying)
}

func (node *Ptr) String() string {
	return fmt.Sprintf("*%s", node.Underlying)
}
