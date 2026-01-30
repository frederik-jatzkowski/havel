package assembly

import (
	"fmt"
	"unsafe"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type call struct {
	fp        bytecode.R
	frameSize uint32
}

func (p *P) AddCall(fp bytecode.R, frameSize uint32, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &call{fp, frameSize})
	p.Positions = append(p.Positions, pos)
}

var _ I = &call{}

func (i *call) ByteCodeLen() int {
	return 1
}

func (i *call) ByteCode(_ int, _ map[string]int) []bytecode.I {
	buf := [4]byte{byte(bytecode.OPCall), byte(i.fp), 0, 0}
	*(*uint16)(unsafe.Pointer(&buf[2])) = uint16(i.frameSize)

	return []bytecode.I{*(*bytecode.I)(unsafe.Pointer(&buf))}
}

func (i *call) String() string {
	return fmt.Sprintf("  call %s %d", i.fp, i.frameSize)
}
