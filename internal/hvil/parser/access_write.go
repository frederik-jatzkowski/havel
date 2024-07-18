package parser

type Write interface{}

type WriteRegister struct {
	Identifier string `"$" @Identifier`
	Type       Type   `":" @@`
}

type WriteVariable struct {
	Identifier string `@Identifier`
}
