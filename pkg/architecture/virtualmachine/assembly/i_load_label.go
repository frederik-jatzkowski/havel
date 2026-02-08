package assembly

import (
	"fmt"
	"unsafe"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type loadLabel struct {
	result bytecode.R
	target string
}

func (p *P) AddLoadLabel(result bytecode.R, label string, pos lexer.Position) {
	p.Instructions = append(p.Instructions, &loadLabel{result, label})
	p.Positions = append(p.Positions, pos, pos, pos)
}

var _ I = &loadLabel{}

func (i *loadLabel) ByteCodeLen() int {
	return 3
}

func (i *loadLabel) ByteCode(_ int, labels map[string]int) []bytecode.I {
	offset := labels[i.target]
	bc := []bytecode.I{*(*bytecode.I)(unsafe.Pointer(&[4]byte{byte(bytecode.OPLit64), byte(i.result), 0, 0})), 0, 0}
	*(*uint64)(unsafe.Pointer(&bc[1])) = uint64(offset)

	return bc
}

func (i *loadLabel) String(_ map[string]int) string {
	return fmt.Sprintf("  load_label %s %s", i.result, i.target)
}
