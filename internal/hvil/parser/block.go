package parser

type BasicBlock struct {
	Identifier   string          `"block":Keyword @Identifier "{"`
	Instructions []Instruction   `@@*`
	JumpTarget   BlockTerminator `"}" "=>" @@ ";"`
}
