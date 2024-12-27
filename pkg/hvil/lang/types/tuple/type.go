package tuple

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
)

type Type struct {
	tool.Node[Type]
	tool.NotImplemented[Type]

	Members []types.Type `parser:"'[' @@ (',' @@)* ']'"`
}

func (t Type) String() string {
	result := "["
	for _, member := range t.Members {
		result += member.String()
	}

	return result + "]"
}

func (t Type) Equals(other types.Type) bool {
	return t.String() == other.String()
}
