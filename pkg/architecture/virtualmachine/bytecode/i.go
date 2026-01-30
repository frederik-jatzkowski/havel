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
	return OP(i & 0xff)
}

func (i I) Regs() (R, R, R) {
	return R(i >> 8 & 0xff), R(i >> 16 & 0xff), R(i >> 24 & 0xff)
}

func (i I) R1Uint16() (R, uint16) {
	return R(i >> 8 & 0xff), uint16(i >> 16 & 0xffff)
}

func (i I) Int16() (int16, int16) {
	return int16(i & 0xffff), int16(i >> 16 & 0xffff)
}
