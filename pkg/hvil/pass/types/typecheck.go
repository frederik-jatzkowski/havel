package names

type TypeCheck[T any] struct {
	TypeCheckPass T `parser:"" json:",omitempty"`
}
