package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type Function struct {
	Pos                    lexer.Position
	pkg                    *Package
	Name                   string       `"func":Keyword @Identifier`
	Head                   FunctionHead `@@`
	Body                   FunctionBody `@@`
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

	for _, declaration := range function.Head.Parameters.Items {
		declaration.GenerateBackLinks(function)
	}

	for _, declaration := range function.Head.ReturnValues.Items {
		declaration.GenerateBackLinks(function)
	}

	for _, declaration := range function.Body.LocalDeclarations.Items {
		declaration.GenerateBackLinks(function)
	}

	for _, block := range function.Body.BasicBlocks {
		block.GenerateBackLinks(function)
	}
}

func (function *Function) ResolveNames(errorsCollector *errors.Collector) {
	function.variableDeclarationMap = make(
		map[string]*FunctionVariableDeclaration,
		len(function.Head.Parameters.Items)+
			len(function.Head.ReturnValues.Items)+
			len(function.Body.LocalDeclarations.Items),
	)

	for _, declaration := range function.Head.Parameters.Items {
		declaration.ResolveNames(errorsCollector)
	}

	for _, declaration := range function.Head.ReturnValues.Items {
		declaration.ResolveNames(errorsCollector)
	}

	for _, declaration := range function.Body.LocalDeclarations.Items {
		declaration.ResolveNames(errorsCollector)
	}

	function.blockMap = make(map[string]*BasicBlock, len(function.Body.BasicBlocks))
	for _, block := range function.Body.BasicBlocks {
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

	for _, block := range function.Body.BasicBlocks {
		block.ResolveNames(errorsCollector)
	}
}
