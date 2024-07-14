package parser

type Literal struct {
}

type FunctionCall struct {
	ReturnValue  string   `@Identifier`
	FunctionName string   `@Identifier`
	Args         []string `@Identifier`
}
