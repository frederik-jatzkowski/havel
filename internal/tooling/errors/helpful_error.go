package errors

import (
	"fmt"
	"io"

	"github.com/alecthomas/participle/v2/lexer"
)

type HelpfulError struct {
	Pos   lexer.Position
	Name  string
	Short string
}

func (err HelpfulError) tryToWriteTo(writer io.Writer) {
	fmt.Fprintf(writer, `
%s at %s: %s
`,
		err.Name,
		err.Pos,
		err.Short,
	)
}
