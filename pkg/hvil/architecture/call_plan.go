package architecture

type CallParam struct {
	BoundTo Register
	RelAddr int
	Bytes   int
}

type CallPlan struct {
	Offset int
	Params []CallParam
}
