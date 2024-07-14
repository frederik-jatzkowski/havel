package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type Function struct {
	Pos                    lexer.Position
	Name                   string       `"func":Keyword @Identifier`
	Head                   FunctionHead `@@`
	Body                   FunctionBody `@@`
	VariableDeclarationMap map[string]*VariableDeclaration
}

func (function *Function) ResolveNames(errorsCollector *errors.Collector) {
	function.VariableDeclarationMap = make(
		map[string]*VariableDeclaration,
		len(function.Head.Parameters.Items)+
			len(function.Head.ReturnValues.Items)+
			len(function.Body.LocalDeclarations.Items),
	)

	for _, declaration := range function.Head.Parameters.Items {
		_, exists := function.VariableDeclarationMap[declaration.Name]
		if exists {
			errorsCollector.Err(
				function.Pos,
				"NameError",
				fmt.Sprintf("the variable %s is redeclared in this function", declaration.Name),
			)
		}

		function.VariableDeclarationMap[declaration.Name] = &declaration
	}

	for _, declaration := range function.Head.ReturnValues.Items {
		_, exists := function.VariableDeclarationMap[declaration.Name]
		if exists {
			errorsCollector.Err(
				function.Pos,
				"NameError",
				fmt.Sprintf("the variable %s is redeclared in this function", declaration.Name),
			)
		}

		function.VariableDeclarationMap[declaration.Name] = &declaration
	}

	for _, declaration := range function.Body.LocalDeclarations.Items {
		_, exists := function.VariableDeclarationMap[declaration.Name]
		if exists {
			errorsCollector.Err(
				function.Pos,
				"NameError",
				fmt.Sprintf("the variable %s is redeclared in this function", declaration.Name),
			)
		}

		function.VariableDeclarationMap[declaration.Name] = &declaration
	}
}

type FunctionHead struct {
	Parameters   CommaSeparatedList[VariableDeclaration] `"(" @@ ")"`
	ReturnValues CommaSeparatedList[VariableDeclaration] `( "=>" "(" @@ ")" )?`
}

type FunctionBody struct {
	LocalDeclarations CommaSeparatedList[VariableDeclaration] `"{" ( "declare":Keyword "(" @@ ")" ";" )?`
	BasicBlocks       []BasicBlock                            `@@+  "}"`
}
