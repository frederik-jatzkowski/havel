package parser

type BlockTerminator interface {
	VisitCLR(visitor Visitor)
}
