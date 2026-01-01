package memory

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type VarWrite struct {
	tool.Node[Write]
	names.NameResolution[struct {
		Decl *stack.Decl
	}]

	Ident string `parser:"@Ident"`
}

func (w *VarWrite) ResolveNames(
	vars names.Scope[*stack.Decl],
	_ names.Scope[*RegWrite],
) (errs []error) {
	decl, err := vars.Find(w.Ident)
	if err != nil {
		return append(errs, w.Wrap(err))
	}

	w.NameResolutionPass.Decl = decl

	return nil
}

func (w *VarWrite) Type() types.Type {
	return w.NameResolutionPass.Decl.Type()
}

func (w *VarWrite) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return w.NameResolutionPass.Decl.Addr(vm)
}
