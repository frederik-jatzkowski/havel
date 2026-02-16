package controlflow

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
)

type Node interface {
	ID() statistics.BlockID
	Successors() []Node
}
