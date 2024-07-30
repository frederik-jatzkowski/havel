package parser

type Read interface {
	VisitLCR(visitor Visitor)
}
