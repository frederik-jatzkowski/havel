package architecture

type MemoryAllocation struct {
	BoundTo Register
	RelAddr int
	Bytes   int
}

type CallPlan struct {
	Offset int
	Params []MemoryAllocation
	Result MemoryAllocation
}
