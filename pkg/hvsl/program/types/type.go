package types

import "github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/names"

type Type interface {
	names.Resolver
	Equals(other Type) bool
}

type TypedObject interface {
	Type() Type
}
