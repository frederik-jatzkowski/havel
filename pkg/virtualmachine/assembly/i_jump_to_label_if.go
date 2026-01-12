package assembly

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type jumpToLabelIf struct {
	condition bytecode.R
	target    string
}

func (p *P) AddJumpToLabelIf(condition bytecode.R, target string, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &jumpToLabelIf{condition, target})
	p.Positions = append(p.Positions, pos)
}

var _ I = &jumpToLabelIf{}

func (i *jumpToLabelIf) ByteCodeLen() int {
	return 1
}

func (i *jumpToLabelIf) ByteCode(index int, labels map[string]int) []bytecode.I {
	offset := labels[i.target] - index
	if offset > math.MaxInt16 || offset < math.MinInt16 {
		panic(fmt.Sprintf("jump target out of range: %d", offset))
	}

	buf := [4]byte{byte(bytecode.OPJumpRelativeIf), byte(i.condition), 0, 0}
	*(*int16)(unsafe.Pointer(&buf[2])) = int16(offset)

	return []bytecode.I{*(*bytecode.I)(unsafe.Pointer(&buf))}
}

func (i *jumpToLabelIf) String() string {
	return fmt.Sprintf("  jump_to_label_if %s %s", i.condition, i.target)
}
