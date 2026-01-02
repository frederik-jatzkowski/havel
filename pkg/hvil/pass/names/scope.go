package names

import (
	"encoding/json"
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type ScopedObject interface {
	tool.NodeLike
	Identifier() string
}

type Scope[T ScopedObject] struct {
	kind  string
	inner *Scope[T]
	defs  map[string]T
}

func NewRootScope[T ScopedObject](kind string) Scope[T] {
	return Scope[T]{
		kind:  kind,
		inner: nil,
		defs:  make(map[string]T),
	}
}

func (s Scope[T]) Define(entry T) error {
	identifier := entry.Identifier()

	if _, err := s.Find(identifier); err == nil {
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

	if s.inner != nil {
		return s.inner.Find(identifier)
	}

	return result, fmt.Errorf("%s '%s' not found", s.kind, identifier)
}

func (s Scope[T]) NewChild() Scope[T] {
	return Scope[T]{
		kind:  s.kind,
		inner: &s,
		defs:  make(map[string]T),
	}
}

func (s Scope[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.defs)
}
