package bytecode

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
)

type P struct {
	Positions    []lexer.Position
	Instructions []I
}

func (p *P) String() string {
	result := ""
	for _, instruction := range p.Instructions {
		result += fmt.Sprintf("%s\n", instruction)
	}

	return result
}
