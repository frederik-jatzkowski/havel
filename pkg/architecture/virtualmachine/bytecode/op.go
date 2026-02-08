package bytecode

import "fmt"

//go:generate go tool go-enum ./op.go

// OP represents a vm instruction set opcode
// ENUM(
//
//		unknown,
//		exit,
//		jump_relative,
//		jump_relative_if,
//		call,
//		lit8,
//		lit16,
//		lit32,
//		lit64,
//		debug_dump,
//		debug_stack_push,
//		debug_stack_pop,
//		alu_add_u8,
//		alu_add_u16,
//		alu_add_u32,
//		alu_add_u64,
//		alu_sub_u8,
//		alu_sub_u16,
//		alu_sub_u32,
//		alu_sub_u64,
//		alu_mul_u8,
//		alu_mul_u16,
//		alu_mul_u32,
//		alu_mul_u64,
//		alu_div_u8,
//		alu_div_u16,
//		alu_div_u32,
//		alu_div_u64,
//		alu_mod_u8,
//		alu_mod_u16,
//		alu_mod_u32,
//		alu_mod_u64,
//		alu_lt_u,
//		alu_eq,
//		alu_move,
//		alloc,
//		free,
//		static_ptr,
//		stack_ptr,
//		store_8,
//		store_16,
//		store_32,
//		store_64,
//		load_8,
//		load_16,
//		load_32,
//		load_64
//	)
type OP byte

func LoadForSize(size int) (OP, error) {
	switch size {
	case 1:
		return OPLoad8, nil
	case 2:
		return OPLoad16, nil
	case 4:
		return OPLoad32, nil
	case 8:
		return OPLoad64, nil
	default:
		return 0, fmt.Errorf("unsupported size %d", size)
	}
}

func StoreForSize(size int) (OP, error) {
	switch size {
	case 1:
		return OPStore8, nil
	case 2:
		return OPStore16, nil
	case 4:
		return OPStore32, nil
	case 8:
		return OPStore64, nil
	default:
		return 0, fmt.Errorf("unsupported size %d", size)
	}
}
