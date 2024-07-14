package parser

type CommaSeparatedList[T any] struct {
	Items []T `(@@ ( "," @@ )* )?`
}
