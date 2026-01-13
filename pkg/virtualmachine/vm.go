package virtualmachine

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type VM struct {
	pc, sp    *int64
	done      bool
	exitCode  int
	registers [32]uint64
	stack     []byte

	stdin          io.Reader
	stdout, stderr io.Writer
}

func New(
	stackSize int,
	stdin io.Reader,
	stdout, stderr io.Writer,
) *VM {
	vm := &VM{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		stack:  make([]byte, stackSize),
	}

	vm.pc = (*int64)(unsafe.Pointer(&vm.registers[0]))
	vm.sp = (*int64)(unsafe.Pointer(&vm.registers[1]))

	return vm
}

func (vm *VM) Execute(p *bytecode.P) error {
	for !vm.done {
		if err := vm.execI(p); err != nil {
			if _, err := fmt.Fprintln(vm.stdout, err.Error()); err != nil {
				panic(err)
			}

			return err
		}
	}

	return nil
}

func (vm *VM) execI(p *bytecode.P) error {
	i := p.Instructions[*vm.pc]
	op := i.OP()
	switch op {
	case bytecode.OPExit:
		r1, _, _ := i.Regs()
		vm.done = true
		vm.exitCode = int(vm.registers[r1])
		*vm.pc++
	case bytecode.OPJumpRelative:
		_, offset := i.Int16()
		*vm.pc += int64(offset)
	case bytecode.OPJumpRelativeIf:
		r1, _, _ := i.Regs()
		if vm.registers[r1]&1 == 0 {
			*vm.pc++
			break
		}

		_, offset := i.Int16()
		*vm.pc += int64(offset)
	case bytecode.OPLit8:
		r1, r2, _ := i.Regs()
		vm.registers[r1] = uint64(r2)
		*vm.pc++
	case bytecode.OPLit16:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = (uint64(r2) << 8) | uint64(r3)
		*vm.pc++
	case bytecode.OPLit32:
		r1, _, _ := i.Regs()
		vm.registers[r1] = uint64(*(*uint32)(unsafe.Pointer(&p.Instructions[*vm.pc+1])))
		*vm.pc += 2
	case bytecode.OPLit64:
		r1, _, _ := i.Regs()
		vm.registers[r1] = *(*uint64)(unsafe.Pointer(&p.Instructions[*vm.pc+1]))
		*vm.pc += 3
	case bytecode.OPDebugDump:
		r1, _, _ := i.Regs()
		if _, err := fmt.Fprintf(vm.stdout, "%s register content: %d\n", p.Positions[*vm.pc], vm.registers[r1]); err != nil {
			panic(err)
		}
		*vm.pc++
	case bytecode.OPAluAddU8:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint8(vm.registers[r2]) + uint8(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluAddU16:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint16(vm.registers[r2]) + uint16(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluAddU32:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint32(vm.registers[r2]) + uint32(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluAddU64:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = vm.registers[r2] + vm.registers[r3]
		*vm.pc++
	case bytecode.OPAluSubU8:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint8(vm.registers[r2]) - uint8(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluSubU16:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint16(vm.registers[r2]) - uint16(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluSubU32:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint32(vm.registers[r2]) - uint32(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluSubU64:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = vm.registers[r2] - vm.registers[r3]
		*vm.pc++
	case bytecode.OPAluMulU8:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint8(vm.registers[r2]) * uint8(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluMulU16:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint16(vm.registers[r2]) * uint16(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluMulU32:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint32(vm.registers[r2]) * uint32(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluMulU64:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = vm.registers[r2] * vm.registers[r3]
		*vm.pc++
	case bytecode.OPAluDivU8:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc])
		}
		vm.registers[r1] = uint64(uint8(vm.registers[r2]) / uint8(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluDivU16:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc])
		}
		vm.registers[r1] = uint64(uint16(vm.registers[r2]) / uint16(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluDivU32:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc])
		}
		vm.registers[r1] = uint64(uint32(vm.registers[r2]) / uint32(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluDivU64:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc])
		}
		vm.registers[r1] = vm.registers[r2] / vm.registers[r3]
		*vm.pc++
	case bytecode.OPAluModU8:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc])
		}
		vm.registers[r1] = uint64(uint8(vm.registers[r2]) % uint8(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluModU16:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc])
		}
		vm.registers[r1] = uint64(uint16(vm.registers[r2]) % uint16(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluModU32:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc])
		}
		vm.registers[r1] = uint64(uint32(vm.registers[r2]) % uint32(vm.registers[r3]))
		*vm.pc++
	case bytecode.OPAluModU64:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[*vm.pc])
		}
		vm.registers[r1] = vm.registers[r2] % vm.registers[r3]
		*vm.pc++
	case bytecode.OPAluLtU:
		r1, r2, r3 := i.Regs()
		res := vm.registers[r2] < vm.registers[r3]
		vm.registers[r1] = uint64(*(*uint8)(unsafe.Pointer(&res)))
		*vm.pc++
	case bytecode.OPAluEq:
		r1, r2, r3 := i.Regs()
		res := vm.registers[r2] == vm.registers[r3]
		vm.registers[r1] = uint64(*(*uint8)(unsafe.Pointer(&res)))
		*vm.pc++
	case bytecode.OPAluMove:
		r1, r2, _ := i.Regs()
		vm.registers[r1] = vm.registers[r2]
		*vm.pc++
	case bytecode.OPStoreStack8:
		r1, _, _ := i.Regs()
		_, offset := i.Uint16()
		vm.stack[*vm.sp+int64(offset)] = uint8(vm.registers[r1])
		*vm.pc++
	case bytecode.OPStoreStack16:
		r1, _, _ := i.Regs()
		_, offset := i.Uint16()
		*(*uint16)(unsafe.Pointer(&vm.stack[*vm.sp+int64(offset)])) = uint16(vm.registers[r1])
		*vm.pc++
	case bytecode.OPStoreStack32:
		r1, _, _ := i.Regs()
		_, offset := i.Uint16()
		*(*uint32)(unsafe.Pointer(&vm.stack[*vm.sp+int64(offset)])) = uint32(vm.registers[r1])
		*vm.pc++
	case bytecode.OPStoreStack64:
		r1, _, _ := i.Regs()
		_, offset := i.Uint16()
		*(*uint64)(unsafe.Pointer(&vm.stack[*vm.sp+int64(offset)])) = vm.registers[r1]
		*vm.pc++
	case bytecode.OPLoadStack8:
		r1, _, _ := i.Regs()
		_, offset := i.Uint16()
		vm.registers[r1] = uint64(vm.stack[*vm.sp+int64(offset)])
		*vm.pc++
	case bytecode.OPLoadStack16:
		r1, _, _ := i.Regs()
		_, offset := i.Uint16()
		vm.registers[r1] = uint64(*(*uint16)(unsafe.Pointer(&vm.stack[*vm.sp+int64(offset)])))
		*vm.pc++
	case bytecode.OPLoadStack32:
		r1, _, _ := i.Regs()
		_, offset := i.Uint16()
		vm.registers[r1] = uint64(*(*uint32)(unsafe.Pointer(&vm.stack[*vm.sp+int64(offset)])))
		*vm.pc++
	case bytecode.OPLoadStack64:
		r1, _, _ := i.Regs()
		_, offset := i.Uint16()
		vm.registers[r1] = *(*uint64)(unsafe.Pointer(&vm.stack[*vm.sp+int64(offset)]))
		*vm.pc++
	default:
		panic(fmt.Sprintf("invalid opcode: %d (%s)", i.OP(), i.OP()))
	}

	return nil
}
