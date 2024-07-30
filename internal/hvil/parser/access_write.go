package parser

type Write interface {
	GenerateBackLinks(*BasicBlock)
	VisitLCR(visitor Visitor)
}
