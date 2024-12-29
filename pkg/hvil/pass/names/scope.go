package names

import (
	"encoding/json"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type ScopedObject interface {
	tool.NodeLike
	Identifier() string
}

type Scope[T ScopedObject] struct {
	kind  string
	inner *Scope[T]
	defs  map[string]*T
}

func NewRootScope[T ScopedObject](kind string) Scope[T] {
	return Scope[T]{
		kind:  kind,
		inner: nil,
		defs:  make(map[string]*T),
	}
}

func (s Scope[T]) Define(entry *T) (err error) {
	identifier := (*entry).Identifier()

	_, exists := s.Find(identifier)
	if exists {
		return (*entry).Errorf("%s '%s' is redeclared", s.kind, identifier)
	}

	s.defs[identifier] = entry

	return nil
}

func (s Scope[T]) DefineAll(entries []T) (errs []error) {
	for i := 0; i < len(entries); i++ {
		err := s.Define(&entries[i])
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (s Scope[T]) Find(identifier string) (*T, bool) {
	result, exists := s.defs[identifier]
	if exists {
		return result, exists
	}

	if s.inner != nil {
		return s.inner.Find(identifier)
	}

	return nil, false
}

func (s Scope[T]) NewChild() Scope[T] {
	return Scope[T]{
		kind:  s.kind,
		inner: &s,
		defs:  make(map[string]*T),
	}
}

func (s Scope[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.defs)
}
