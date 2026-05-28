package types

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type StructMember struct {
	tool.Node[StructMember]

	Name string   `parser:"@Ident"`
	Typ  TypeName `parser:"@@"`
}

func (node *StructMember) Identifier() string {
	return node.Name
}

func (node *StructMember) ResolveNames(ctx context.Context) error {
	return node.Typ.ResolveNames(ctx)
}

func (node *StructMember) Type() Type {
	return node.Typ.Type()
}

type Struct struct {
	tool.Node[Struct]
	names.NameResolution[struct {
		Done    bool
		Decl    *Decl
		Members scope.Scope[*StructMember]
	}]

	Members tool.List[StructMember] `parser:"'struct':Keyword '{' @@ '}'"`
}

func (node *Struct) ResolveNames(ctx context.Context) error {
	if node.NameResolutionPass.Done {
		return nil
	}

	node.NameResolutionPass.Done = true

	decl, err := contexttool.CurrentFromContext[*Decl](ctx)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Decl = decl

	node.NameResolutionPass.Members = scope.NewRoot[*StructMember](names.KindStructMember)

	for i := range node.Members.Items {
		member := &node.Members.Items[i]

		if err := node.NameResolutionPass.Members.Define(member); err != nil {
			return err
		}

		if err := member.ResolveNames(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (node *Struct) Equals(other Type) bool {
	otherStruct, ok := other.(*Struct)
	if !ok {
		return false
	}

	return node.NameResolutionPass.Decl == otherStruct.NameResolutionPass.Decl
}

func (node *Struct) String() string {
	return node.NameResolutionPass.Decl.Identifier()
}
