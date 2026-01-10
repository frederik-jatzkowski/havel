package registeralloc

import "github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"

type Value interface {
	Register() architecture.Register
}
