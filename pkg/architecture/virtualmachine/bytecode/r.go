package bytecode

import (
	"fmt"
)

type R byte

const (
	PC R = 0
	SP R = 1
)

func (r R) String() string {
	return r.RegisterName()
}

func (r R) RegisterName() string {
	switch r {
	case PC:
		return "pc"
	case SP:
		return "sp"
	default:
		return fmt.Sprintf("r%d", r)
	}
}
