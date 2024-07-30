package parser

type Operation interface {
	VisitCLR(visitor Visitor)
}
