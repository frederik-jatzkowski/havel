package assembly

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type I interface {
	fmt.Stringer
	ByteCode() []bytecode.I
}
