package instruction

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type MemoryRead interface {
	tool.NodeLike
	names.Resolver
	statistics.Calculator
	registeralloc.Value
	codegen.VirtualMachine
	Type() types.Type
}

type MemoryWrite interface {
	names.Resolver
	statistics.Calculator
	registeralloc.Value
	codegen.VirtualMachine
	Type() types.Type
}
