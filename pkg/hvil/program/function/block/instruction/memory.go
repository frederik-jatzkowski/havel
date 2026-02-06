package instruction

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
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
