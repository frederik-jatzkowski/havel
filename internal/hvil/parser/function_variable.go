package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type FunctionVariableDeclaration struct {
	Name         string `parser:"@Identifier"`
	DeclaredType Type   `parser:"':' @@"`
	Pos          lexer.Position
	Tokens       []lexer.Token
}

func (declaration *FunctionVariableDeclaration) VisitCLR(visitor Visitor) {
	visitor.VisitFunctionVariableDeclaration(declaration)
}

func (declaration *FunctionVariableDeclaration) Type() Type {
	return declaration.DeclaredType
}
