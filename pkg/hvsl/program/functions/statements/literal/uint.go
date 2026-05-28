package literal

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Uint struct {
	tool.Node[Uint]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Value   uint64 `parser:"'uint' '(' @Number"`
	BitSize uint8  `parser:"',' @Number ')'"`
}

func (node *Uint) ResolveNames(ctx context.Context) error {
	return nil
}

func (node *Uint) ResolveTypes(ctx context.Context) error {
	switch node.BitSize {
	case 8:
		node.TypeCheckPass.Type = types.BuiltinU8
	case 16:
		node.TypeCheckPass.Type = types.BuiltinU16
	case 32:
		node.TypeCheckPass.Type = types.BuiltinU32
	case 64:
		node.TypeCheckPass.Type = types.BuiltinU64
	default:
		return node.Errorf("unsupported bit size for uint: %d", node.BitSize)
	}

	return nil
}

func (node *Uint) Type() types.Type {
	return node.TypeCheckPass.Type
}
