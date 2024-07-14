package parser

import (
	"fmt"
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/frederik-jatzkowski/havel/internal/hvil/token"
)

var parser = participle.MustBuild[Package](
	participle.Lexer(token.Tokenizer),
	participle.Elide("Whitespace", "Comment"),
	participle.UseLookahead(1),
	participle.Union[Type](PrimitiveType{}, TupleType{}),
	participle.Union[JumpTarget](Return{}, Jump{}, ConditionalJump{}),
	participle.Union[ReadAccess](RegisterReadAccess{}, VariableReadAccess{}),
	participle.Union[WriteAccess](RegisterWriteAccess{}, VariableWriteAccess{}),
)

func Parse(fileName string, reader io.Reader) (Package, error) {
	pkg, err := parser.Parse(fileName, reader)
	if err != nil {
		return Package{}, err
	}

	if pkg == nil {
		return Package{}, fmt.Errorf("parser did not return a parsed value")
	}

	return *pkg, err
}
