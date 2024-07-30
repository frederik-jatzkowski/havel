package parser

import "github.com/alecthomas/participle/v2/lexer"

type Read interface {
	VisitCLR(visitor Visitor)
	Type() Type
	Position() lexer.Position
}
