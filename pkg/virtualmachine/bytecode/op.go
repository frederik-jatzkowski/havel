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
//		alu_add_u_8
//	)
type OP byte
