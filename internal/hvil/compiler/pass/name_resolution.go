package pass

type NameResolution[T any] struct {
	NameResolutionPass T `parser:"" json:",omitempty"`
}
