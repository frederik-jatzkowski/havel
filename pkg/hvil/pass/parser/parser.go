package parser

import (
	"fmt"
	"io"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/terminator"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime/alu"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime/debug"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime/literal"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime/local"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/token"

	"github.com/alecthomas/participle/v2"
)

var parser = participle.MustBuild[program.Program](
	participle.Lexer(token.Tokenizer),
	participle.Elide("Whitespace", "Comment"),
	participle.UseLookahead(1),
	participle.Union[types.Type](&types.ScalarType{}, &types.TupleType{}, &types.FunctionType{}),
	participle.Union[block.Terminator](&terminator.Return{}, &terminator.Jump{}, &terminator.Conditional{}),
	participle.Union[memory.Write](&memory.RegWrite{}, &memory.VarWrite{}),
	participle.Union[memory.Read](&memory.RegRead{}, &memory.VarRead{}),
	participle.Union[instruction.Operation](&literal.Scalar{}, &alu.Call{}, &local.Call{}, &debug.Call{}),
	participle.Union[alu.Operation](
		&alu.Move{},
		&alu.AddU{},
		&alu.SubU{},
		&alu.EQ{},
		&alu.LtU{},
	),
	participle.Union[debug.Operation](&debug.Dump{}),
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
