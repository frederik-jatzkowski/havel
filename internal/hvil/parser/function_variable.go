package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type FunctionVariableDeclaration struct {
	Name   string `@Identifier`
	Type   Type   `":" @@`
	Pos    lexer.Position
	Tokens []lexer.Token
}

func (declaration *FunctionVariableDeclaration) VisitLCR(visitor Visitor) {
	visitor.VisitFunctionVariableDeclaration(declaration)
}
