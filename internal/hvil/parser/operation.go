package parser

type Operation interface {
	GenerateBackLinks(*BasicBlock)
	VisitLCR(visitor Visitor)
}
