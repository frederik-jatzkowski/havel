package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type FunctionVariableDeclaration struct {
	Name     string `@Identifier`
	Type     Type   `":" @@`
	Pos      lexer.Position
	Tokens   []lexer.Token
	function *Function
}

func (declaration *FunctionVariableDeclaration) GenerateBackLinks(function *Function) {
	declaration.function = function
}

func (declaration *FunctionVariableDeclaration) ResolveNames(errorsCollector *errors.Collector) {
	_, exists := declaration.function.variableDeclarationMap[declaration.Name]
	if exists {
		errorsCollector.Err(
			declaration.Pos,
			"NameError",
			fmt.Sprintf("the variable %s is redeclared in this function", declaration.Name),
		)
	}

	declaration.function.variableDeclarationMap[declaration.Name] = declaration
}
