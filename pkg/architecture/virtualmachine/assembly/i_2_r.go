package assembly

import (
	"fmt"
	"unsafe"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type i2R struct {
	op     bytecode.OP
	r1, r2 bytecode.R
}

func (p *P) AddI2R(op bytecode.OP, r1, r2 bytecode.R, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &i2R{op, r1, r2})
	p.Positions = append(p.Positions, pos)
}

var _ I = &i2R{}

func (i *i2R) ByteCodeLen() int {
	return 1
}

func (i *i2R) ByteCode(_ int, _ map[string]int) []bytecode.I {
	return []bytecode.I{
		*(*bytecode.I)(unsafe.Pointer(&[4]byte{byte(i.op), byte(i.r1), byte(i.r2), 0})),
	}
}

func (i *i2R) String(_ map[string]int) string {
	return fmt.Sprintf("  %s %s %s", i.op, i.r1, i.r2)
}
