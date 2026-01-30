package assembly

import (
	"fmt"
	"unsafe"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type i3R struct {
	op         bytecode.OP
	r1, r2, r3 bytecode.R
}

func (p *P) AddI3R(op bytecode.OP, r1, r2, r3 bytecode.R, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &i3R{op, r1, r2, r3})
	p.Positions = append(p.Positions, pos)
}

var _ I = &i3R{}

func (i *i3R) ByteCodeLen() int {
	return 1
}

func (i *i3R) ByteCode(_ int, _ map[string]int) []bytecode.I {
	return []bytecode.I{
		*(*bytecode.I)(unsafe.Pointer(&[4]byte{byte(i.op), byte(i.r1), byte(i.r2), byte(i.r3)})),
	}
}

func (i *i3R) String() string {
	return fmt.Sprintf("  %s %s %s %s", i.op, i.r1, i.r2, i.r3)
}
