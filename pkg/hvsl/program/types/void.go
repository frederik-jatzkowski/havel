package types

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Void struct {
	tool.Node[Struct]
	Void string `parser:"'void'"`
}

func (node *Void) ResolveNames(ctx context.Context) error {
	return nil
}

func (node *Void) Equals(other Type) bool {
	_, ok := other.(*Void)

	return ok
}
