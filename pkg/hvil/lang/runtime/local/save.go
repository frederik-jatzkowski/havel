package local

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
)

type SavedMemory struct {
	Register architecture.Register
	RelAddr  int
	Bytes    int
}
