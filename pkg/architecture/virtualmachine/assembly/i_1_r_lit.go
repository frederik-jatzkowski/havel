package assembly

import (
	"fmt"
	"unsafe"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type i1RLit struct {
	op  bytecode.OP
	r1  bytecode.R
	lit uint16
}

func (p *P) AddI1RLit(op bytecode.OP, r1 bytecode.R, lit uint16, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &i1RLit{op, r1, lit})
	p.Positions = append(p.Positions, pos)
}

var _ I = &i1RLit{}

func (i *i1RLit) ByteCodeLen() int {
	return 1
}

func (i *i1RLit) ByteCode(_ int, _ map[string]int) []bytecode.I {
	buf := [4]byte{byte(i.op), byte(i.r1), 0, 0}
	*(*uint16)(unsafe.Pointer(&buf[2])) = i.lit

	return []bytecode.I{*(*bytecode.I)(unsafe.Pointer(&buf))}
}

func (i *i1RLit) String(_ map[string]int) string {
	return fmt.Sprintf("  %s %s %d", i.op, i.r1, i.lit)
}
