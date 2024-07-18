package parser

type Read interface{}

type ReadRegister struct {
	Identifier string `"$" @Identifier`
}

type ReadVariable struct {
	Identifier string `@Identifier`
}
