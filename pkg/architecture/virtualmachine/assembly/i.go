package assembly

import (
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type I interface {
	String(labels map[string]int) string
	ByteCodeLen() int
	ByteCode(i int, labels map[string]int) []bytecode.I
}
