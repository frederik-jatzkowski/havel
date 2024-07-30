package parser

type Write interface {
	VisitLCR(visitor Visitor)
}
