package types

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type TupleType struct {
	tool.Node[TupleType]
	tool.NotImplemented[TupleType]

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
	return node.Bytes() == other.Bytes()
}

func (node *TupleType) Bytes() (size int) {
	for _, member := range node.Members {
		size += member.Bytes()
	}

	return size
}
