package literal

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Boolean struct {
	tool.Node[Boolean]

	True  bool `parser:"('true' | "`
	False bool `parser:"'false' )"`
}

func (node *Boolean) ResolveNames(ctx context.Context) error {
	return nil
}

func (node *Boolean) ResolveTypes(ctx context.Context) error {
	return nil
}

func (node *Boolean) Type() types.Type {
	return types.BuiltinBool
}
