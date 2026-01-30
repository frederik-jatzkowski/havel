package assembly

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
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
	labelMap := make(map[string]int)
	i := 0
	for _, instr := range p.Instructions {

		l, ok := instr.(*label)
		if !ok {
			i += instr.ByteCodeLen()
			continue
		}

		labelMap[l.name] = i
	}

	byteCode := &bytecode.P{
		Positions: p.Positions,
	}

	i = 0
	for _, instr := range p.Instructions {
		byteCode.Instructions = append(byteCode.Instructions, instr.ByteCode(i, labelMap)...)
		i += instr.ByteCodeLen()
	}

	return byteCode, nil
}
