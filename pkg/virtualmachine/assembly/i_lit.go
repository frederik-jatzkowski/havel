package assembly

import (
	"fmt"
	"unsafe"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type lit struct {
	target bytecode.R
	size   int
	value  uint64
}

func (p *P) AddLit(target bytecode.R, size int, value uint64, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &lit{target, size, value})

	switch size {
	case 1, 2:
		p.Positions = append(p.Positions, pos)
	case 4:
		p.Positions = append(p.Positions, pos, pos)
	case 8:
		p.Positions = append(p.Positions, pos, pos, pos)
	}
}

var _ I = &lit{}

func (i *lit) ByteCode() []bytecode.I {
	switch i.size {
	case 1:
		return []bytecode.I{*(*bytecode.I)(unsafe.Pointer(&[4]byte{byte(bytecode.OPLit1), byte(i.target), uint8(i.value), 0}))}
	case 2:
		buf := [4]byte{byte(bytecode.OPLit2), byte(i.target), 0, 0}
		*(*uint16)(unsafe.Pointer(&buf[3])) = uint16(i.value)

		return []bytecode.I{*(*bytecode.I)(unsafe.Pointer(&buf))}
	case 4:
		bc := []bytecode.I{*(*bytecode.I)(unsafe.Pointer(&[4]byte{byte(bytecode.OPLit4), byte(i.target), 0, 0})), 0}
		*(*uint32)(unsafe.Pointer(&bc[1])) = uint32(i.value)

		return bc
	case 8:
		bc := []bytecode.I{*(*bytecode.I)(unsafe.Pointer(&[4]byte{byte(bytecode.OPLit8), byte(i.target), 0, 0})), 0, 0}
		*(*uint64)(unsafe.Pointer(&bc[1])) = i.value

		return bc
	default:
		panic(fmt.Sprintf("invalid size: %d", i.size))
	}
}

func (i *lit) String() string {
	return fmt.Sprintf("  lit_%d %s %d", i.size, i.target, i.value)
}
