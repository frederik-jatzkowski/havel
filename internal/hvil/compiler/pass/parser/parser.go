package parser

import (
	"fmt"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/function"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/function/block/terminator"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/register"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types/scalar"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types/tuple"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/variable"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/runtime/alu"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/runtime/debug"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/pass/token"
	"io"

	"github.com/alecthomas/participle/v2"
)

var parser = participle.MustBuild[program.Program](
	participle.Lexer(token.Tokenizer),
	participle.Elide("Whitespace", "Comment"),
	participle.UseLookahead(1),
	participle.Union[types.Type](scalar.Type{}, tuple.Type{}, function.Type{}),
	participle.Union[block.Terminator](&terminator.Return{}, &terminator.Jump{}, &terminator.Conditional{}),
	participle.Union[memory.Read](&register.Read{}, &variable.Read{}),
	participle.Union[memory.Write](&register.Write{}, &variable.Write{}),
	participle.Union[block.Operation](&scalar.Literal{}, &alu.Operation{}, &function.Call{}, &debug.Prefix{}),
	participle.Union[debug.Op](&debug.PrintU32{}),
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
