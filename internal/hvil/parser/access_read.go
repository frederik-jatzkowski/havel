package parser

type Read interface {
	GenerateBackLinks(*BasicBlock)
	VisitLCR(visitor Visitor)
}
