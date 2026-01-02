package names

type Resolver interface {
	ResolveNames() error
}

type NameResolution[T any] struct {
	NameResolutionPass T `parser:"" json:",omitempty"`
}
