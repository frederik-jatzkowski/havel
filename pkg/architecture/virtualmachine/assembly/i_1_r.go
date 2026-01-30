package assembly

import (
	"fmt"
	"unsafe"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type i1R struct {
	op bytecode.OP
	r1 bytecode.R
}

func (p *P) AddI1R(op bytecode.OP, r1 bytecode.R, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &i1R{op, r1})
	p.Positions = append(p.Positions, pos)
}

var _ I = &i1R{}

func (i *i1R) ByteCodeLen() int {
	return 1
}

func (i *i1R) ByteCode(_ int, _ map[string]int) []bytecode.I {
	return []bytecode.I{
		*(*bytecode.I)(unsafe.Pointer(&[4]byte{byte(i.op), byte(i.r1), 0, 0})),
	}
}

func (i *i1R) String() string {
	return fmt.Sprintf("  %s %s", i.op, i.r1)
}
