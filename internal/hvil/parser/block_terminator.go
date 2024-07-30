package parser

type BlockTerminator interface {
	VisitLCR(visitor Visitor)
}
