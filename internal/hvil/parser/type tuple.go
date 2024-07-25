package parser

type TupleType struct {
	Members []Type `"[" @@ ("," @@)* "]"`
}
