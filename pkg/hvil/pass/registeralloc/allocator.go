package registeralloc

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
)

type Allocator interface {
	NewScope() Scope
}

type allocator struct {
	arch architecture.Architecture
}

func NewAllocator(arch architecture.Architecture) Allocator {
	return &allocator{
		arch: arch,
	}
}

func (a *allocator) NewScope() Scope {
	return &scope{
		generalPurposeRegisters: a.arch.GeneralPurposeRegisters(),
		scratchRegisters:        a.arch.ScratchRegisters(),
	}
}
