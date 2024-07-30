package parser

import "encoding/json"

type CommaSeparatedList[T any] struct {
	Items []T `parser:"(@@ ( ',' @@ )* )?"`
}

func (list CommaSeparatedList[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.Items)
}
