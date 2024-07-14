package parser

type BasicBlock struct {
	Identifier string     `"block":Keyword @Identifier "{" "}"`
	JumpTarget JumpTarget `"=>" @@ ";"`
}

type JumpTarget interface{}

type Return struct {
	Token string `@"return":Keyword`
}

type Jump struct {
	Target string `@Identifier`
}

type ConditionalJump struct {
	Condition ReadAccess `"if":Keyword @@`
	True      string     `"then":Keyword @Identifier`
	False     string     `"else":Keyword @Identifier`
}
