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

func (function *Function) VisitCLR(visitor Visitor) {
	visitor.SetCurrentFunction(function)
	visitor.VisitFunction(function)

	for _, declaration := range function.Parameters.Items {
		declaration.VisitCLR(visitor)
	}

	if function.ReturnValue != nil {
		function.ReturnValue.VisitCLR(visitor)
	}

	for _, declaration := range function.LocalDeclarations.Items {
		declaration.VisitCLR(visitor)
	}

	for _, block := range function.BasicBlocks {
		block.VisitCLR(visitor)
	}
}
