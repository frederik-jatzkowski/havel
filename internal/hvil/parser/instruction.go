package parser

import "github.com/alecthomas/participle/v2/lexer"

type Instruction struct {
	Pos       lexer.Position
	block     *BasicBlock
	Result    *Write    `(@@ "=")?`
	Operation Operation `@@ ";"`
}

func (instr *Instruction) GenerateBackLinks(block *BasicBlock) {
	instr.block = block
}
