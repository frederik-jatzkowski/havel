package parser

type Operation interface {
}

type PrimitiveLiteral struct {
	Value uint64 `@BitLiteral`
}

type AluOperation struct {
	Name string `"alu" "." @Identifier`
	Arg1 Read   `"(" @@`
	Arg2 *Read  `("," @@)? ")"`
}

type LocalCall struct {
	Name string                   `"local" "." @Identifier`
	Args CommaSeparatedList[Read] `"(" @@ ")"`
}

type DebugOperation struct {
	Name string `"debug" "." @Identifier`
	Arg1 Read   `"(" @@ ")"`
}
