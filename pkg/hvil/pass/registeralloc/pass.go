package registeralloc

type RegisterAllocation[T any] struct {
	RegisterAllocationPass T `parser:"" json:",omitempty"`
}
