package registeralloc

import "github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"

type Range struct {
	Start, End statistics.InstructionID
}
