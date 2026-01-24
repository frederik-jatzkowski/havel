package controlflow

import "github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"

type LiveRange struct {
	Start, End *statistics.InstructionID
}

func ComputeLiveRanges(
	entry Node,
	usages, defs map[statistics.BlockID][]statistics.InstructionID,
) map[statistics.BlockID][]LiveRange {
	blocks := discoverBlocks(entry)
	liveIn, liveOut := computeBlockLevelLiveness(blocks, usages, defs)

	result := make(map[statistics.BlockID][]LiveRange, len(blocks))
	for _, block := range blocks {
		blockResult := make([]LiveRange, 0)
		if liveIn[block.ID()] {
			blockResult = append(blockResult, LiveRange{})
		}

		accesses := CalculateOrderedAccessList(defs[block.ID()], usages[block.ID()])
		for _, access := range accesses {
			switch access.Kind {
			case AccessKindWRITE:
				blockResult = append(blockResult, LiveRange{Start: &access.ID})
			case AccessKindREAD:
				blockResult[len(blockResult)-1].End = &access.ID
			}
		}

		if liveOut[block.ID()] {
			blockResult[len(blockResult)-1].End = nil
		}

		result[block.ID()] = blockResult
	}

	return result
}

func computeBlockLevelLiveness(
	blocks []Node,
	usages map[statistics.BlockID][]statistics.InstructionID,
	defs map[statistics.BlockID][]statistics.InstructionID,
) (liveIn, liveOut map[statistics.BlockID]bool) {
	liveIn = make(map[statistics.BlockID]bool)
	liveOut = make(map[statistics.BlockID]bool)

	for {
		changed := false
		// Iterate in reverse order for faster convergence in backward analysis
		for i := len(blocks) - 1; i >= 0; i-- {
			b := blocks[i]
			id := b.ID()

			// OUT[n] = Union of IN[s] for all successors s
			newLiveOut := false
			for _, succ := range b.Successors() {
				if liveIn[succ.ID()] {
					newLiveOut = true
					break
				}
			}

			// IN[n] = Use[n] OR (OUT[n] AND NOT Def[n])
			hasUse := len(usages[id]) > 0
			hasDef := len(defs[id]) > 0

			// To be precise: Is there a use BEFORE a def in this block?
			usedBeforeDef := false
			if hasUse {
				if !hasDef || usages[id][0] <= defs[id][0] {
					usedBeforeDef = true
				}
			}

			newLiveIn := usedBeforeDef || (newLiveOut && !hasDef)

			if liveIn[id] != newLiveIn || liveOut[id] != newLiveOut {
				liveIn[id] = newLiveIn
				liveOut[id] = newLiveOut
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	return liveIn, liveOut
}

func MustBeSavedAt(ranges []LiveRange, id statistics.InstructionID) bool {
	for _, liveRange := range ranges {
		wasLiveBefore := liveRange.Start == nil || *liveRange.Start < id
		willBeLiveAfter := liveRange.End == nil || *liveRange.End > id

		if wasLiveBefore && willBeLiveAfter {
			return true
		}
	}

	return false
}
