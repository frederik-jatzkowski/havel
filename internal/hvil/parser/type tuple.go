package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type TupleType struct {
	Members []Type `parser:"'[' @@ (',' @@)* ']'"`
	Tokens  []lexer.Token
}

func (t TupleType) String() string {
	result := "["
	for _, member := range t.Members {
		result += member.String()
	}

	return result + "]"
}

func (t TupleType) Equals(other Type) bool {
	return false
}
