package parser

import (
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type Operation interface {
	GenerateBackLinks(*BasicBlock)
	ResolveNames(*errors.Collector)
}
