package bytecode

//go:generate go tool go-enum ./op.go

// OP represents a vm instruction set opcode
// ENUM(
//
//		unknown,
//		exit,
//		lit8,
//		lit16,
//		lit32,
//		lit64,
//		debug_dump,
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
//		store_stack_8,
//		store_stack_16,
//		store_stack_32,
//		store_stack_64,
//		load_stack_8,
//		load_stack_16,
//		load_stack_32,
//		load_stack_64
//	)
type OP byte
