package tool

import "encoding/json"

type List[T any] struct {
	Items []T `parser:"(@@ ( ',' @@ )* )?"`
}

func (list List[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.Items)
}
