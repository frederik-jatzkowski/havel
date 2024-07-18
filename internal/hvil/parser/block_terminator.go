package parser

type BlockTerminator interface{}

type Return struct {
	Token string `@"return":Keyword`
}

type Jump struct {
	Target string `@Identifier`
}

type ConditionalJump struct {
	Condition Read   `"if":Keyword @@`
	True      string `"then":Keyword @Identifier`
	False     string `"else":Keyword @Identifier`
}
