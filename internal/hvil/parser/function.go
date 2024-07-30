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
	variableDeclarationMap map[string]*FunctionVariableDeclaration
	blockMap               map[string]*BasicBlock
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
