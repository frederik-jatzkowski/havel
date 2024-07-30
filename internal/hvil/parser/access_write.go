package parser

type Write interface {
	VisitCLR(visitor Visitor)
}
