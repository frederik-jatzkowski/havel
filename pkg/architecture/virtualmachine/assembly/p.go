package assembly

import (
	"bytes"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type P struct {
	Positions    []lexer.Position
	StaticData   []S
	Instructions []I
}

func NewP() *P {
	return &P{}
}

func (p *P) String() string {
	var result strings.Builder

	result.WriteString(".static\n")

	for _, data := range p.StaticData {
		result.WriteString(data.String())
		result.WriteString("\n")
	}

	result.WriteString("\n.code\n")

	for _, instruction := range p.Instructions {
		result.WriteString(instruction.String())
		result.WriteString("\n")
	}

	return result.String()
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

	buf := bytes.NewBuffer(nil)
	for _, s := range p.StaticData {
		if _, err := s.WriteTo(buf); err != nil {
			return nil, err
		}
	}

	byteCode.StaticData = buf.Bytes()

	i = 0
	for _, instr := range p.Instructions {
		byteCode.Instructions = append(byteCode.Instructions, instr.ByteCode(i, labelMap)...)
		i += instr.ByteCodeLen()
	}

	return byteCode, nil
}
