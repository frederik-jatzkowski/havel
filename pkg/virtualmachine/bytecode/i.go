package bytecode

import (
	"fmt"
	"unsafe"
)

type I uint32

func (i I) String() string {
	components := *(*[4]byte)(unsafe.Pointer(&i))
	return fmt.Sprintf("%08b", components)
}

func (i I) OP() OP {
	return *(*OP)(unsafe.Pointer(&i))
}

func (i I) Regs() (R, R, R) {
	arr := *(*[4]R)(unsafe.Pointer(&i))
	return arr[1], arr[2], arr[3]
}

func (i I) Uint16() (uint16, uint16) {
	arr := *(*[2]uint16)(unsafe.Pointer(&i))
	return arr[0], arr[1]
}
