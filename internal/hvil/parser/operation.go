package parser

type Operation interface {
	VisitLCR(visitor Visitor)
}
