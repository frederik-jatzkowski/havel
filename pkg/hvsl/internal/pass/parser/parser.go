package parser

import (
	"fmt"
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/token"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements/alu"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements/call"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements/debug"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements/literal"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements/mem"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
)

var parser = participle.MustBuild[program.Program](
	participle.Lexer(token.Tokenizer),
	participle.Elide("Whitespace", "Comment"),
	participle.UseLookahead(1),
	participle.Union[program.Member](
		&functions.Function{},
		&types.Decl{},
	),
	participle.Union[functions.Member](
		&statements.Statement{},
		&controlflow.If{},
		&controlflow.For{},
	),
	participle.Union[statements.Content](
		&statements.Return{},
		&statements.Let{},
		&statements.Evaluation{},
		&debug.Call{},
		&alu.Call{},
		&literal.Call{},
		&mem.Call{},
		&call.Call{},
	),
	participle.Union[statements.Expression](
		&alu.Call{},
		&literal.Call{},
		&mem.Call{},
		&mem.Ident{},
		&call.Call{},
	),
	participle.Union[debug.Operation](
		&debug.Dump{},
	),
	participle.Union[literal.Operation](
		&literal.Boolean{},
		&literal.Uint{},
	),
	participle.Union[alu.Operation](
		&alu.UnOp{},
		&alu.BinOp{},
	),
	participle.Union[mem.Operation](
		&mem.Ptr{},
		&mem.GEP{},
		&mem.Load{},
		&mem.Store{},
		&mem.Alloc{},
	),
	participle.Union[call.Operation](
		&call.Local{},
	),
	participle.Union[types.Type](
		&types.Struct{},
	),
)

func Parse(fileName string, reader io.Reader) (program.Program, error) {
	prog, err := parser.Parse(fileName, reader)
	if err != nil {
		return program.Program{}, err
	}

	if prog == nil {
		return program.Program{}, fmt.Errorf("lang did not return a parsed value")
	}

	return *prog, err
}
