package registeralloc

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
)

type Range struct {
	Start, End statistics.InstructionID
}
