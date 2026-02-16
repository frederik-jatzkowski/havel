package typecheck

type TypeCheck[T any] struct {
	TypeCheckPass T `parser:"" json:",omitempty"`
}
