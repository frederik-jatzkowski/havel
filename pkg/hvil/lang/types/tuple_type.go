package types

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type TupleType struct {
	tool.Node[TupleType]
	tool.NotImplemented[TupleType]

	Members []Type `parser:"'[' @@ (',' @@)* ']'"`
}

func (t TupleType) String() string {
	result := "["
	for i, member := range t.Members {
		result += member.String()
		if i < len(t.Members)-1 {
			result += ", "
		}
	}

	return result + "]"
}

func (t TupleType) CanBeAssigned(other Type) bool {
	return t.BitSize() == other.BitSize()
}

func (t TupleType) BitSize() (size int) {
	for _, member := range t.Members {
		size += member.BitSize()
	}

	return size
}
