package bytecode

//go:generate go tool go-enum ./op.go

// OP represents a vm instruction set opcode
// ENUM(
//
//		unknown,
//		exit,
//		lit_1,
//		lit_2,
//		lit_4,
//		lit_8,
//		debug_dump,
//		alu_add_u_1,
//		alu_add_u_2,
//		alu_add_u_4,
//		alu_add_u_8,
//		alu_sub_u_1,
//		alu_sub_u_2,
//		alu_sub_u_4,
//		alu_sub_u_8,
//		alu_mul_u_1,
//		alu_mul_u_2,
//		alu_mul_u_4,
//		alu_mul_u_8,
//		alu_div_u_1,
//		alu_div_u_2,
//		alu_div_u_4,
//		alu_div_u_8,
//		alu_mod_u_1,
//		alu_mod_u_2,
//		alu_mod_u_4,
//		alu_mod_u_8,
//		alu_lt_u,
//		alu_eq,
//		alu_move
//	)
type OP byte
