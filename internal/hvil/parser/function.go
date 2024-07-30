package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Function struct {
	Name                   string                                           `"func":Keyword @Identifier`
	Parameters             CommaSeparatedList[*FunctionVariableDeclaration] `"(" @@ ")"`
	ReturnValue            *FunctionVariableDeclaration                     `( "=>" "(" @@ ")" )?`
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

	if function.ReturnValue != nil {
		function.ReturnValue.GenerateBackLinks(function)
	}

	for _, declaration := range function.LocalDeclarations.Items {
		declaration.GenerateBackLinks(function)
	}

	for _, block := range function.BasicBlocks {
		block.GenerateBackLinks(function)
	}
}

func (function *Function) VisitLCR(visitor Visitor) {
	visitor.SetCurrentFunction(function)
	visitor.VisitFunction(function)

	for _, declaration := range function.Parameters.Items {
		declaration.VisitLCR(visitor)
	}

	if function.ReturnValue != nil {
		function.ReturnValue.VisitLCR(visitor)
	}

	for _, declaration := range function.LocalDeclarations.Items {
		declaration.VisitLCR(visitor)
	}

	for _, block := range function.BasicBlocks {
		block.VisitLCR(visitor)
	}
}
