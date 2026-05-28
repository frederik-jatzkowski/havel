package functions

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Member interface {
	tool.NodeLike
	names.Resolver
	typecheck.Resolver
}
