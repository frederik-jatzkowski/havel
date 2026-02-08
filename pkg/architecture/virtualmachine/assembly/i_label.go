package assembly

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type label struct {
	name string
}

func (p *P) AddLabel(name string, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &label{name})
}

var _ I = &label{}

func (i *label) ByteCodeLen() int {
	return 0
}

func (i *label) ByteCode(_ int, _ map[string]int) []bytecode.I {
	return nil
}

func (i *label) String(labels map[string]int) string {
	return fmt.Sprintf("%s (pc %d):", i.name, labels[i.name])
}
