package virtualmachine

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
)

type VM struct {
	done     bool
	exitCode int

	pc, sp *int64

	registers [32]uint64

	stackSize int

	heap *Heap

	stdin          io.Reader
	stdout, stderr io.Writer
}

func New(
	stackSize int,
	stdin io.Reader,
	stdout, stderr io.Writer,
) *VM {
	vm := &VM{
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
		stackSize: stackSize,
	}

	vm.pc = (*int64)(unsafe.Pointer(&vm.registers[0]))
	vm.sp = (*int64)(unsafe.Pointer(&vm.registers[1]))

	vm.heap = &Heap{}
	vm.heap.data = [][]byte{make([]byte, 0), make([]byte, uint64(vm.stackSize))}

	return vm
}

func (vm *VM) Execute(p *bytecode.P) (err error) {
	vm.done = false
	vm.exitCode = 0

	*vm.pc = 0
	*vm.sp = int64(NewFatPtr(1, 0).ToUint64())

	if len(vm.heap.data) > 1 {
		vm.heap.data = [][]byte{vm.heap.data[0]}
		vm.heap.free = nil
	}

	vm.heap.data = [][]byte{make([]byte, len(p.StaticData)), make([]byte, uint64(vm.stackSize))}
	vm.heap.free = nil

	copy(vm.heap.data[0], p.StaticData)

	defer func() {
		r := recover()
		if r != nil {
			fmt.Fprintln(vm.stdout, r)
			if e, ok := r.(error); ok {
				err = e
			}
		}
	}()

	for !vm.done {
		vm.execI(p)
	}

	return err
}

func (vm *VM) execI(p *bytecode.P) {
	i := p.Instructions[*vm.pc]
	op := i.OP()
	switch op {
	case bytecode.OPExit:
		vm.execOPExit(p, i)
	case bytecode.OPJumpRelative:
		vm.execOPJumpRelative(p, i)
	case bytecode.OPJumpRelativeIf:
		vm.execOPJumpRelativeIf(p, i)
	case bytecode.OPCall:
		vm.execOPCall(p, i)
	case bytecode.OPLit8:
		vm.execOPLit8(p, i)
	case bytecode.OPLit16:
		vm.execOPLit16(p, i)
	case bytecode.OPLit32:
		vm.execOPLit32(p, i)
	case bytecode.OPLit64:
		vm.execOPLit64(p, i)
	case bytecode.OPDebugDump:
		vm.execOPDebugDump(p, i)
	case bytecode.OPAluAddU8:
		vm.execOPAluAddU8(p, i)
	case bytecode.OPAluAddU16:
		vm.execOPAluAddU16(p, i)
	case bytecode.OPAluAddU32:
		vm.execOPAluAddU32(p, i)
	case bytecode.OPAluAddU64:
		vm.execOPAluAddU64(p, i)
	case bytecode.OPAluSubU8:
		vm.execOPAluSubU8(p, i)
	case bytecode.OPAluSubU16:
		vm.execOPAluSubU16(p, i)
	case bytecode.OPAluSubU32:
		vm.execOPAluSubU32(p, i)
	case bytecode.OPAluSubU64:
		vm.execOPAluSubU64(p, i)
	case bytecode.OPAluMulU8:
		vm.execOPAluMulU8(p, i)
	case bytecode.OPAluMulU16:
		vm.execOPAluMulU16(p, i)
	case bytecode.OPAluMulU32:
		vm.execOPAluMulU32(p, i)
	case bytecode.OPAluMulU64:
		vm.execOPAluMulU64(p, i)
	case bytecode.OPAluDivU8:
		vm.execOPAluDivU8(p, i)
	case bytecode.OPAluDivU16:
		vm.execOPAluDivU16(p, i)
	case bytecode.OPAluDivU32:
		vm.execOPAluDivU32(p, i)
	case bytecode.OPAluDivU64:
		vm.execOPAluDivU64(p, i)
	case bytecode.OPAluModU8:
		vm.execOPAluModU8(p, i)
	case bytecode.OPAluModU16:
		vm.execOPAluModU16(p, i)
	case bytecode.OPAluModU32:
		vm.execOPAluModU32(p, i)
	case bytecode.OPAluModU64:
		vm.execOPAluModU64(p, i)
	case bytecode.OPAluLtU:
		vm.execOPAluLtU(p, i)
	case bytecode.OPAluEq:
		vm.execOPAluEq(p, i)
	case bytecode.OPAluMove:
		vm.execOPAluMove(p, i)
	case bytecode.OPAlloc:
		vm.execOPAlloc(p, i)
	case bytecode.OPFree:
		vm.execOPFree(p, i)
	case bytecode.OPStaticPtr:
		vm.execOPStaticPtr(p, i)
	case bytecode.OPStackPtr:
		vm.execOPStackPtr(p, i)
	case bytecode.OPStore8:
		vm.execOPStore8(p, i)
	case bytecode.OPStore16:
		vm.execOPStore16(p, i)
	case bytecode.OPStore32:
		vm.execOPStore32(p, i)
	case bytecode.OPStore64:
		vm.execOPStore64(p, i)
	case bytecode.OPLoad8:
		vm.execOPLoad8(p, i)
	case bytecode.OPLoad16:
		vm.execOPLoad16(p, i)
	case bytecode.OPLoad32:
		vm.execOPLoad32(p, i)
	case bytecode.OPLoad64:
		vm.execOPLoad64(p, i)
	default:
		panic(fmt.Sprintf("invalid opcode: %d (%s)", i.OP(), i.OP()))
	}
}

//go:inline
func (vm *VM) execOPExit(p *bytecode.P, i bytecode.I) {
	r1, _, _ := i.Regs()
	vm.done = true
	vm.exitCode = int(vm.registers[r1])
	*vm.pc++
}

//go:inline
func (vm *VM) execOPJumpRelative(p *bytecode.P, i bytecode.I) {
	offset := int16(i >> 16 & 0xffff)
	*vm.pc += int64(offset)
}

//go:inline
func (vm *VM) execOPJumpRelativeIf(p *bytecode.P, i bytecode.I) {
	r1, _, _ := i.Regs()
	if vm.registers[r1]&1 == 0 {
		*vm.pc++
		return
	}

	offset := int16(i >> 16 & 0xffff)
	*vm.pc += int64(offset)
}

//go:inline
func (vm *VM) execOPCall(p *bytecode.P, i bytecode.I) {
	fp, frameSize := i.R1Uint16()

	// advance stack pointer
	vm.heap.Store64(NewFatPtr(1, uint32(*vm.sp+int64(frameSize)+8)), uint64(*vm.sp))
	*vm.sp += int64(frameSize)

	// prepare return address
	vm.heap.Store64(NewFatPtr(1, uint32(*vm.sp)), uint64(*vm.pc))

	newPC := int64(vm.registers[fp])

	*vm.pc = newPC
}

//go:inline
func (vm *VM) execOPLit8(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	vm.registers[r1] = uint64(r2)
	*vm.pc++
}

//go:inline
func (vm *VM) execOPLit16(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = (uint64(r2) << 8) | uint64(r3)
	*vm.pc++
}

//go:inline
func (vm *VM) execOPLit32(p *bytecode.P, i bytecode.I) {
	r1, _, _ := i.Regs()
	vm.registers[r1] = uint64(*(*uint32)(unsafe.Pointer(&p.Instructions[*vm.pc+1])))
	*vm.pc += 2
}

//go:inline
func (vm *VM) execOPLit64(p *bytecode.P, i bytecode.I) {
	r1, _, _ := i.Regs()
	vm.registers[r1] = *(*uint64)(unsafe.Pointer(&p.Instructions[*vm.pc+1]))
	*vm.pc += 3
}

//go:inline
func (vm *VM) execOPDebugDump(p *bytecode.P, i bytecode.I) {
	r1, _, _ := i.Regs()
	if _, err := fmt.Fprintf(vm.stdout, "%s register content: %d\n", p.Positions[*vm.pc], vm.registers[r1]); err != nil {
		panic(err)
	}
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluAddU8(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = uint64(uint8(vm.registers[r2]) + uint8(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluAddU16(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = uint64(uint16(vm.registers[r2]) + uint16(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluAddU32(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = uint64(uint32(vm.registers[r2]) + uint32(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluAddU64(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = vm.registers[r2] + vm.registers[r3]
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluSubU8(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = uint64(uint8(vm.registers[r2]) - uint8(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluSubU16(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = uint64(uint16(vm.registers[r2]) - uint16(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluSubU32(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = uint64(uint32(vm.registers[r2]) - uint32(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluSubU64(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = vm.registers[r2] - vm.registers[r3]
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluMulU8(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = uint64(uint8(vm.registers[r2]) * uint8(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluMulU16(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = uint64(uint16(vm.registers[r2]) * uint16(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluMulU32(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = uint64(uint32(vm.registers[r2]) * uint32(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluMulU64(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	vm.registers[r1] = vm.registers[r2] * vm.registers[r3]
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluDivU8(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	if vm.registers[r3] == 0 {
		panic(fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc]))
	}
	vm.registers[r1] = uint64(uint8(vm.registers[r2]) / uint8(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluDivU16(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	if vm.registers[r3] == 0 {
		panic(fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc]))
	}
	vm.registers[r1] = uint64(uint16(vm.registers[r2]) / uint16(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluDivU32(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	if vm.registers[r3] == 0 {
		panic(fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc]))
	}
	vm.registers[r1] = uint64(uint32(vm.registers[r2]) / uint32(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluDivU64(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	if vm.registers[r3] == 0 {
		panic(fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc]))
	}
	vm.registers[r1] = vm.registers[r2] / vm.registers[r3]
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluModU8(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	if vm.registers[r3] == 0 {
		panic(fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc]))
	}
	vm.registers[r1] = uint64(uint8(vm.registers[r2]) % uint8(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluModU16(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	if vm.registers[r3] == 0 {
		panic(fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc]))
	}
	vm.registers[r1] = uint64(uint16(vm.registers[r2]) % uint16(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluModU32(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	if vm.registers[r3] == 0 {
		panic(fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc]))
	}
	vm.registers[r1] = uint64(uint32(vm.registers[r2]) % uint32(vm.registers[r3]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluModU64(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	if vm.registers[r3] == 0 {
		panic(fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc]))
	}
	vm.registers[r1] = vm.registers[r2] % vm.registers[r3]
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluLtU(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	res := vm.registers[r2] < vm.registers[r3]
	vm.registers[r1] = uint64(*(*uint8)(unsafe.Pointer(&res)))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluEq(p *bytecode.P, i bytecode.I) {
	r1, r2, r3 := i.Regs()
	res := vm.registers[r2] == vm.registers[r3]
	vm.registers[r1] = uint64(*(*uint8)(unsafe.Pointer(&res)))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAluMove(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	vm.registers[r1] = vm.registers[r2]
	*vm.pc++
}

//go:inline
func (vm *VM) execOPAlloc(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	ptr := vm.heap.Alloc(vm.registers[r2])
	vm.registers[r1] = ptr.ToUint64()
	*vm.pc++
}

//go:inline
func (vm *VM) execOPFree(p *bytecode.P, i bytecode.I) {
	r1, _, _ := i.Regs()
	ptr := NewFatPtrFromUint64(vm.registers[r1])
	vm.heap.Free(ptr)
	*vm.pc++
}

//go:inline
func (vm *VM) execOPStaticPtr(p *bytecode.P, i bytecode.I) {
	r1, offset := i.R1Uint16()
	ptr := NewFatPtr(0, uint32(offset)+uint32(*vm.sp))
	vm.registers[r1] = ptr.ToUint64()
	*vm.pc++
}

//go:inline
func (vm *VM) execOPStackPtr(p *bytecode.P, i bytecode.I) {
	r1, offset := i.R1Uint16()
	ptr := NewFatPtr(1, uint32(offset)+uint32(*vm.sp))
	vm.registers[r1] = ptr.ToUint64()
	*vm.pc++
}

//go:inline
func (vm *VM) execOPStore8(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	ptr := NewFatPtrFromUint64(vm.registers[r1])
	vm.heap.Store8(ptr, uint8(vm.registers[r2]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPStore16(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	ptr := NewFatPtrFromUint64(vm.registers[r1])
	vm.heap.Store16(ptr, uint16(vm.registers[r2]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPStore32(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	ptr := NewFatPtrFromUint64(vm.registers[r1])
	vm.heap.Store32(ptr, uint32(vm.registers[r2]))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPStore64(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	ptr := NewFatPtrFromUint64(vm.registers[r1])
	vm.heap.Store64(ptr, vm.registers[r2])
	*vm.pc++
}

//go:inline
func (vm *VM) execOPLoad8(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	ptr := NewFatPtrFromUint64(vm.registers[r2])
	vm.registers[r1] = uint64(vm.heap.Load8(ptr))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPLoad16(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	ptr := NewFatPtrFromUint64(vm.registers[r2])
	vm.registers[r1] = uint64(vm.heap.Load16(ptr))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPLoad32(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	ptr := NewFatPtrFromUint64(vm.registers[r2])
	vm.registers[r1] = uint64(vm.heap.Load32(ptr))
	*vm.pc++
}

//go:inline
func (vm *VM) execOPLoad64(p *bytecode.P, i bytecode.I) {
	r1, r2, _ := i.Regs()
	ptr := NewFatPtrFromUint64(vm.registers[r2])
	vm.registers[r1] = vm.heap.Load64(ptr)
	*vm.pc++
}
