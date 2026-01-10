package virtualmachine

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type VM struct {
	pc        int
	done      bool
	exitCode  int
	registers [256]uint64

	stdin          io.Reader
	stdout, stderr io.Writer
}

func New(
	stackSize int,
	stdin io.Reader,
	stdout, stderr io.Writer,
) *VM {
	return &VM{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
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
	i := p.Instructions[vm.pc]
	switch i.OP() {
	case bytecode.OPExit:
		r1, _, _ := i.Regs()
		vm.done = true
		vm.exitCode = int(vm.registers[r1])
		vm.pc++
	case bytecode.OPLit1:
		r1, r2, _ := i.Regs()
		vm.registers[r1] = uint64(r2)
		vm.pc++
	case bytecode.OPLit2:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = (uint64(r2) << 8) | uint64(r3)
		vm.pc++
	case bytecode.OPLit4:
		r1, _, _ := i.Regs()
		vm.registers[r1] = uint64(*(*uint32)(unsafe.Pointer(&p.Instructions[vm.pc+1])))
		vm.pc += 2
	case bytecode.OPLit8:
		r1, _, _ := i.Regs()
		vm.registers[r1] = *(*uint64)(unsafe.Pointer(&p.Instructions[vm.pc+1]))
		vm.pc += 3
	case bytecode.OPDebugDump:
		r1, _, _ := i.Regs()
		if _, err := fmt.Fprintf(vm.stdout, "%s register content: %d\n", p.Positions[vm.pc], vm.registers[r1]); err != nil {
			panic(err)
		}
		vm.pc++
	case bytecode.OPAluAddU1:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint8(vm.registers[r2]) + uint8(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluAddU2:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint16(vm.registers[r2]) + uint16(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluAddU4:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint32(vm.registers[r2]) + uint32(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluAddU8:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = vm.registers[r2] + vm.registers[r3]
		vm.pc++
	case bytecode.OPAluSubU1:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint8(vm.registers[r2]) - uint8(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluSubU2:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint16(vm.registers[r2]) - uint16(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluSubU4:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint32(vm.registers[r2]) - uint32(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluSubU8:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = vm.registers[r2] - vm.registers[r3]
		vm.pc++
	case bytecode.OPAluMulU1:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint8(vm.registers[r2]) * uint8(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluMulU2:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint16(vm.registers[r2]) * uint16(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluMulU4:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = uint64(uint32(vm.registers[r2]) * uint32(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluMulU8:
		r1, r2, r3 := i.Regs()
		vm.registers[r1] = vm.registers[r2] * vm.registers[r3]
		vm.pc++
	case bytecode.OPAluDivU1:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[vm.pc])
		}
		vm.registers[r1] = uint64(uint8(vm.registers[r2]) / uint8(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluDivU2:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[vm.pc])
		}
		vm.registers[r1] = uint64(uint16(vm.registers[r2]) / uint16(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluDivU4:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[vm.pc])
		}
		vm.registers[r1] = uint64(uint32(vm.registers[r2]) / uint32(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluDivU8:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[vm.pc])
		}
		vm.registers[r1] = vm.registers[r2] / vm.registers[r3]
		vm.pc++
	case bytecode.OPAluModU1:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[vm.pc])
		}
		vm.registers[r1] = uint64(uint8(vm.registers[r2]) % uint8(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluModU2:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[vm.pc])
		}
		vm.registers[r1] = uint64(uint16(vm.registers[r2]) % uint16(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluModU4:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[vm.pc])
		}
		vm.registers[r1] = uint64(uint32(vm.registers[r2]) % uint32(vm.registers[r3]))
		vm.pc++
	case bytecode.OPAluModU8:
		r1, r2, r3 := i.Regs()
		if vm.registers[r3] == 0 {
			return fmt.Errorf("%s: division by zero\n", p.Positions[vm.pc])
		}
		vm.registers[r1] = vm.registers[r2] % vm.registers[r3]
		vm.pc++
	case bytecode.OPAluLtU:
		r1, r2, r3 := i.Regs()
		res := vm.registers[r2] < vm.registers[r3]
		vm.registers[r1] = uint64(*(*uint8)(unsafe.Pointer(&res)))
		vm.pc++
	case bytecode.OPAluEq:
		r1, r2, r3 := i.Regs()
		res := vm.registers[r2] == vm.registers[r3]
		vm.registers[r1] = uint64(*(*uint8)(unsafe.Pointer(&res)))
		vm.pc++
	case bytecode.OPAluMove:
		r1, r2, _ := i.Regs()
		vm.registers[r1] = vm.registers[r2]
		vm.pc++
	default:
		panic(fmt.Sprintf("invalid opcode: %d (%s)", i.OP(), i.OP()))
	}

	return nil
}
