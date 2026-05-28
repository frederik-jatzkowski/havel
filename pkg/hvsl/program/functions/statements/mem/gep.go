package mem

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type GEPChainIdentifier struct {
	tool.Node[GEPChainIdentifier]

	Name string `parser:"@Ident"`
}

type GEP struct {
	tool.Node[GEP]
	typecheck.TypeCheck[struct {
		Result types.Type
	}]

	Ptr   statements.Expression `parser:"'gep' '(' @@"`
	Chain []GEPChainIdentifier  `parser:"( '.' @@ )+ ')'"`
}

func (node *GEP) ResolveNames(ctx context.Context) error {
	return node.Ptr.ResolveNames(ctx)
}

func (node *GEP) ResolveTypes(ctx context.Context) error {
	if err := node.Ptr.ResolveTypes(ctx); err != nil {
		return err
	}

	ptr, ok := node.Ptr.Type().(*types.Ptr)
	if !ok {
		return node.Errorf("expected pointer, got %s", node.Ptr.Type())
	}

	node.TypeCheckPass.Result = ptr.Underlying
	for _, item := range node.Chain {
		prev := node.TypeCheckPass.Result

		structDecl, ok := prev.(*types.Struct)
		if !ok {
			return item.Errorf("%s is not a struct", prev)
		}

		member, err := structDecl.NameResolutionPass.Members.Find(item.Name)
		if err != nil {
			return node.Wrap(err)
		}

		node.TypeCheckPass.Result = member.Type()
	}

	node.TypeCheckPass.Result = &types.Ptr{Underlying: node.TypeCheckPass.Result}

	return nil
}

func (node *GEP) Type() types.Type {
	return node.TypeCheckPass.Result
}
