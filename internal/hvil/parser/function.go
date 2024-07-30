package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Function struct {
	Name                   string                                           `parser:"'func':Keyword @Identifier"`
	Parameters             CommaSeparatedList[*FunctionVariableDeclaration] `parser:"'(' @@ ')'"`
	ReturnValue            *FunctionVariableDeclaration                     `parser:"( '=>' '(' @@ ')' )?"`
	LocalDeclarations      CommaSeparatedList[*FunctionVariableDeclaration] `parser:"'{' ( 'declare':Keyword '(' @@ ')' ';' )?"`
	BasicBlocks            []*BasicBlock                                    `parser:"@@+  '}'"`
	Pos                    lexer.Position
	Tokens                 []lexer.Token
	VariableDeclarationMap map[string]*FunctionVariableDeclaration
	BlockMap               map[string]*BasicBlock `parser:"" json:"-"`
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
