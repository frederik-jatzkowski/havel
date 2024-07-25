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
	participle.Union[Type](PrimitiveType{}),
	participle.Union[BlockTerminator](&Return{}, &Jump{}, &ConditionalJump{}),
	participle.Union[Read](&ReadRegister{}, &ReadVariable{}),
	participle.Union[Write](&WriteRegister{}, &WriteVariable{}),
	participle.Union[Operation](&PrimitiveLiteral{}, &AluOperation{}, &LocalCall{}, &DebugOperation{}),
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
