package parser

type WriteAccess interface{}

type RegisterWriteAccess struct {
	Identifier string `"$" @Identifier`
	Type       Type   `":" @@`
}

type VariableWriteAccess struct {
	Identifier string `@Identifier`
}
