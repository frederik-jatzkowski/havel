package literal

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"math/bits"
)

type Scalar struct {
	tool.Node[Scalar]

	Value uint64 `parser:"@BitLiteral"`
}

func (l *Scalar) ResolveNames(vars names.Scope[memory.VarDecl], regs names.Scope[memory.RegWrite]) (errs []error) {
	return nil
}

func (l *Scalar) ResolveTypes(target types.Type) (errs []error) {
	_, ok := target.(types.ScalarType)
	if !ok {
		return append(errs, l.Errorf("cannot assign scalar literal to %s", target))
	}

	requiredBitSize := bits.Len64(l.Value)
	expectedBitSize := target.BitSize()
	if requiredBitSize > expectedBitSize {
		errs = append(errs, l.Errorf("cannot assign scalar literal %d to %s: value too big", l.Value, target))
	}

	return errs
}
