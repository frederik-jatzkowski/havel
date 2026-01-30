package assembly

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type jumpToLabel struct {
	target string
}

func (p *P) AddJumpToLabel(target string, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &jumpToLabel{target})
	p.Positions = append(p.Positions, pos)
}

var _ I = &jumpToLabel{}

func (i *jumpToLabel) ByteCodeLen() int {
	return 1
}

func (i *jumpToLabel) ByteCode(index int, labels map[string]int) []bytecode.I {
	offset := labels[i.target] - index
	if offset > math.MaxInt16 || offset < math.MinInt16 {
		panic(fmt.Sprintf("jump target out of range: %d", offset))
	}

	buf := [4]byte{byte(bytecode.OPJumpRelative), 0, 0, 0}
	*(*int16)(unsafe.Pointer(&buf[2])) = int16(offset)

	return []bytecode.I{*(*bytecode.I)(unsafe.Pointer(&buf))}
}

func (i *jumpToLabel) String() string {
	return fmt.Sprintf("  jump_to_label %s", i.target)
}
