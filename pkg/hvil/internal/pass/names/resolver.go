package names

import "context"

type Resolver interface {
	ResolveNames(ctx context.Context) error
}
