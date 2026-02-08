package assembly

import (
	"fmt"
	"unsafe"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type i0R struct {
	op bytecode.OP
}

func (p *P) AddI0R(op bytecode.OP, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &i0R{op})
	p.Positions = append(p.Positions, pos)
}

var _ I = &i0R{}

func (i *i0R) ByteCodeLen() int {
	return 1
}

func (i *i0R) ByteCode(_ int, _ map[string]int) []bytecode.I {
	return []bytecode.I{
		*(*bytecode.I)(unsafe.Pointer(&[4]byte{byte(i.op), 0, 0, 0})),
	}
}

func (i *i0R) String(_ map[string]int) string {
	return fmt.Sprintf("  %s", i.op)
}
