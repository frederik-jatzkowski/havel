package parser

type VariableDeclaration struct {
	Name string `@Identifier`
	Type Type   `":" @@`
}
