package parser

type Read interface {
	VisitCLR(visitor Visitor)
}
