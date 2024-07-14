package parser

type Type interface{}

type PrimitiveType struct {
	BitSize uint8 `@BitSize`
}

type TupleType struct {
	Members []Type `"[" @@ ("," @@)* "]"`
}
