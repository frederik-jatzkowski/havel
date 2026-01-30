package parser

import (
	"fmt"
	"io"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/token"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction/alu"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction/call"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction/debug"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction/literal"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction/mem"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/terminator"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"

	"github.com/alecthomas/participle/v2"
)

var parser = participle.MustBuild[program.Program](
	participle.Lexer(token.Tokenizer),
	participle.Elide("Whitespace", "Comment"),
	participle.UseLookahead(1),
	participle.Union[function.Block](&block.Block{}),
	participle.Union[types.Type](
		&types.Scalar{},
		&types.Ref{},
		&types.Function{},
	),
	participle.Union[block.Terminator](
		&terminator.Return{},
		&terminator.Jump{},
		&terminator.Conditional{},
	),
	participle.Union[instruction.MemoryWrite](
		&instruction.RegWrite{},
		&instruction.VarWrite{},
	),
	participle.Union[instruction.MemoryRead](
		&instruction.RegRead{},
		&instruction.VarRead{},
	),
	participle.Union[instruction.Operation](
		&literal.Scalar{},
		&alu.Call{},
		&call.Call{},
		&debug.Call{},
		&mem.Call{},
	),
	participle.Union[alu.Operation](
		&alu.Move{},
		&alu.AddU{},
		&alu.SubU{},
		&alu.MulU{},
		&alu.DivU{},
		&alu.ModU{},
		&alu.EQ{},
		&alu.LtU{},
	),
	participle.Union[mem.Operation](
		&mem.Alloc{},
		&mem.Free{},
		&mem.Store{},
		&mem.Load{},
		&mem.Ptr{},
	),
	participle.Union[call.Operation](
		&call.Local{},
		&call.Ptr{},
		&call.Dyn{},
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

func String() string {
	return parser.String()
}
