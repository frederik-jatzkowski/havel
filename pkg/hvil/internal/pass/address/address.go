package address

type Resolution[T any] struct {
	AddressResolutionPass T `parser:"" json:",omitempty"`
}
