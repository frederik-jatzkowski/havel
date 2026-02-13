package scope

import (
	"encoding/json"
	"fmt"
	"iter"
	"maps"
)

type Scope[T Object] struct {
	kind  string
	outer *Scope[T]
	defs  map[string]T
}

func NewRoot[T Object](kind fmt.Stringer) Scope[T] {
	return Scope[T]{
		kind: kind.String(),
		defs: make(map[string]T),
	}
}

func (s Scope[T]) Child() Scope[T] {
	return Scope[T]{
		kind:  s.kind,
		outer: &s,
		defs:  make(map[string]T),
	}
}

func (s Scope[T]) Define(entry T) error {
	identifier := entry.Identifier()

	_, exists := s.defs[identifier]
	if exists {
		return entry.Errorf("%s '%s' is redeclared", s.kind, identifier)
	}

	s.defs[identifier] = entry

	return nil
}

func (s Scope[T]) Find(identifier string) (T, error) {
	result, exists := s.defs[identifier]
	if exists {
		return result, nil
	}

	if s.outer != nil {
		return s.outer.Find(identifier)
	}

	return result, fmt.Errorf("%s '%s' not found", s.kind, identifier)
}

// All returns all elements in the current and all parent scopes.
func (s Scope[T]) All() iter.Seq[T] {
	if s.outer == nil {
		return maps.Values(s.defs)
	}

	return func(yield func(T) bool) {
		for _, iterator := range []iter.Seq[T]{maps.Values(s.defs), s.outer.All()} {
			for each := range iterator {
				if !yield(each) {
					return
				}
			}
		}
	}
}

func (s Scope[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.defs)
}
