package controlflow

import (
	"slices"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
)

type Access struct {
	Kind AccessKind
	ID   statistics.InstructionID
}

func CalculateOrderedAccessList(defs, usages []statistics.InstructionID) []Access {
	accesses := make([]Access, 0, len(defs)+len(usages))

	for _, id := range defs {
		accesses = append(accesses, Access{Kind: AccessKindWRITE, ID: id})
	}

	for _, id := range usages {
		accesses = append(accesses, Access{Kind: AccessKindREAD, ID: id})
	}

	slices.SortFunc(accesses, func(a, b Access) int {
		if a.ID == b.ID {
			return int(a.Kind) - int(b.Kind)
		}

		return int(a.ID) - int(b.ID)
	})

	return accesses
}
