package names

type NameResolution[T any] struct {
	NameResolutionPass T `parser:"" json:",omitempty"`
}
