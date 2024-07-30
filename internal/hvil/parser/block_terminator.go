package parser

type BlockTerminator interface {
	GenerateBackLinks(*BasicBlock)
	VisitLCR(visitor Visitor)
}
