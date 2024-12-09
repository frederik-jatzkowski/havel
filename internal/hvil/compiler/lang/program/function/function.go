package function

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type Function struct {
	tool.Node[Function]

	Name        string                 `parser:"'func':Keyword @Ident"`
	Parameters  tool.List[*stack.Decl] `parser:"'(' @@ ')'"`
	ReturnValue *stack.Decl            `parser:"( '=>' '(' @@ ')' )?"`
	LocalDecls  tool.List[*stack.Decl] `parser:"'{' ( 'declare':Keyword '(' @@ ')' ';' )?"`
	BasicBlocks []*block.Block         `parser:"@@+  '}'"`
}

func (function *Function) Type() (result Type) {
	for _, param := range function.Parameters.Items {
		result.Parameters.Items = append(result.Parameters.Items, param.DeclaredType)
	}

	if function.ReturnValue != nil {
		result.ReturnValue = function.ReturnValue.DeclaredType
	}

	return result
}
