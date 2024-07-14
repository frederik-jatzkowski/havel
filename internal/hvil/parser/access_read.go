package parser

type ReadAccess interface{}

type RegisterReadAccess struct {
	Identifier string `"$" @Identifier`
}

type VariableReadAccess struct {
	Identifier string `@Identifier`
}
