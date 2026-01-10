package assembly

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type P struct {
	Positions    []lexer.Position
	Instructions []I
}

func NewP() *P {
	return &P{}
}

func (p *P) String() string {
	result := ""
	for _, instruction := range p.Instructions {
		result += fmt.Sprintf("%s\n", instruction.String())
	}

	return result
}

func (p *P) Assemble() (*bytecode.P, error) {
	byteCode := &bytecode.P{
		Positions: p.Positions,
	}
	for _, instr := range p.Instructions {
		byteCode.Instructions = append(byteCode.Instructions, instr.ByteCode()...)
	}

	return byteCode, nil
}
