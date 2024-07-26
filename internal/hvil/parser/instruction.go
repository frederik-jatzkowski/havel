package parser

import "github.com/alecthomas/participle/v2/lexer"

type Instruction struct {
	Result    *Write    `(@@ "=")?`
	Operation Operation `@@ ";"`
	Pos       lexer.Position
	block     *BasicBlock
	Tokens    []lexer.Token
}

func (instr *Instruction) GenerateBackLinks(block *BasicBlock) {
	instr.block = block
}
