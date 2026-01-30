package assembly

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type I interface {
	fmt.Stringer
	ByteCodeLen() int
	ByteCode(i int, labels map[string]int) []bytecode.I
}
