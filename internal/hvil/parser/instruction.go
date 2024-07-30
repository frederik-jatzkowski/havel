package parser

import "github.com/alecthomas/participle/v2/lexer"

type Instruction struct {
	Result    *Write    `parser:"(@@ '=')?"`
	Operation Operation `parser:"@@ ';'"`
	Pos       lexer.Position
	Tokens    []lexer.Token
}

func (instr *Instruction) VisitCLR(visitor Visitor) {
	visitor.VisitInstruction(instr)

	if instr.Result != nil {
		(*instr.Result).VisitCLR(visitor)
	}

	instr.Operation.VisitCLR(visitor)
}
