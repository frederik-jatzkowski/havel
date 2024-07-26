package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type Function struct {
	Name                   string                                           `"func":Keyword @Identifier`
	Parameters             CommaSeparatedList[*FunctionVariableDeclaration] `"(" @@ ")"`
	ReturnValues           CommaSeparatedList[*FunctionVariableDeclaration] `( "=>" "(" @@ ")" )?`
	LocalDeclarations      CommaSeparatedList[*FunctionVariableDeclaration] `"{" ( "declare":Keyword "(" @@ ")" ";" )?`
	BasicBlocks            []*BasicBlock                                    `@@+  "}"`
	Pos                    lexer.Position
	Tokens                 []lexer.Token
	pkg                    *Package
	variableDeclarationMap map[string]*FunctionVariableDeclaration
	blockMap               map[string]*BasicBlock
}

type FunctionHead struct {
	Parameters   CommaSeparatedList[*FunctionVariableDeclaration] `"(" @@ ")"`
	ReturnValues CommaSeparatedList[*FunctionVariableDeclaration] `( "=>" "(" @@ ")" )?`
}

type FunctionBody struct {
	LocalDeclarations CommaSeparatedList[*FunctionVariableDeclaration] `"{" ( "declare":Keyword "(" @@ ")" ";" )?`
	BasicBlocks       []*BasicBlock                                    `@@+  "}"`
}

func (function *Function) GenerateBackLinks(pkg *Package) {
	function.pkg = pkg

	for _, declaration := range function.Parameters.Items {
		declaration.GenerateBackLinks(function)
	}

	for _, declaration := range function.ReturnValues.Items {
		declaration.GenerateBackLinks(function)
	}

	for _, declaration := range function.LocalDeclarations.Items {
		declaration.GenerateBackLinks(function)
	}

	for _, block := range function.BasicBlocks {
		block.GenerateBackLinks(function)
	}
}

func (function *Function) ResolveNames(errorsCollector *errors.Collector) {
	function.variableDeclarationMap = make(
		map[string]*FunctionVariableDeclaration,
		len(function.Parameters.Items)+
			len(function.ReturnValues.Items)+
			len(function.LocalDeclarations.Items),
	)

	for _, declaration := range function.Parameters.Items {
		declaration.ResolveNames(errorsCollector)
	}

	for _, declaration := range function.ReturnValues.Items {
		declaration.ResolveNames(errorsCollector)
	}

	for _, declaration := range function.LocalDeclarations.Items {
		declaration.ResolveNames(errorsCollector)
	}

	function.blockMap = make(map[string]*BasicBlock, len(function.BasicBlocks))
	for _, block := range function.BasicBlocks {
		_, exists := function.blockMap[block.Identifier]
		if exists {
			errorsCollector.Err(
				block.Pos,
				"NameError",
				fmt.Sprintf("the basic block %s is redeclared in this function", block.Identifier),
			)
		}

		function.blockMap[block.Identifier] = block
	}

	for _, block := range function.BasicBlocks {
		block.ResolveNames(errorsCollector)
	}
}
