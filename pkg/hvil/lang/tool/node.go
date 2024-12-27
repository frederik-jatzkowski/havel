package tool

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
)

type Node[TKind any] struct {
	Pos    lexer.Position `parser:"" json:"-"`
	EndPos lexer.Position `parser:"" json:"-"`
	Tokens []lexer.Token  `parser:"" json:"-"`

	Kind Kind[TKind] `parser:""` // allows serialization of the token kind name
}

func (n Node[TKind]) String() (s string) {
	for _, t := range n.Tokens {
		s += t.String()
	}

	return s
}

func (n Node[TKind]) Position() lexer.Position {
	return n.Pos
}

func (n Node[TKind]) Errorf(format string, a ...any) error {
	return fmt.Errorf("%s: "+format, append([]any{n.Position()}, a...)...)
}
