package controlflow

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
)

type Node interface {
	ID() statistics.BlockID
	Successors() []Node
}
