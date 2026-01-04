package types

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type TupleType struct {
	tool.Node[TupleType]

	Members []Type `parser:"'[' @@ (',' @@)* ']'"`
}

func (node *TupleType) String() string {
	result := "["
	for i, member := range node.Members {
		result += member.String()
		if i < len(node.Members)-1 {
			result += ", "
		}
	}

	return result + "]"
}

func (node *TupleType) CanBeAssigned(other Type) bool {
	return node.Equals(other)
}

func (node *TupleType) Equals(other Type) bool {
	return node.EqualsDetailed(other) == nil
}

func (node *TupleType) EqualsDetailed(other Type) error {
	otherTuple, ok := other.(*TupleType)
	if !ok {
		return fmt.Errorf("%s is not a tuple type", other)
	}

	if len(otherTuple.Members) != len(node.Members) {
		return fmt.Errorf("tuple length mismatch: expected %d, got %d", len(node.Members), len(otherTuple.Members))
	}

	for i, member := range node.Members {
		otherMember := otherTuple.Members[i]
		if err := member.EqualsDetailed(otherMember); err != nil {
			return fmt.Errorf("tuple member %d mismatch: expected %s, got %s", i, member, otherMember)
		}
	}

	return nil
}

func (node *TupleType) Bytes() (size int) {
	for _, member := range node.Members {
		size += member.Bytes()
	}

	return size
}
