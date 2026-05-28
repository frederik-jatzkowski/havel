package alu

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

//go:generate go tool go-enum ./symbols.go

// BinOpSymbol represents binary operation symbols.
// ENUM(add, sub, mul, div, mod, eq, leq, lt, geq, gt)
type BinOpSymbol byte

func (node *BinOpSymbol) Parse(lex *lexer.PeekingLexer) error {
	token := lex.Peek()

	parsed, err := ParseBinOpSymbol(token.Value)
	if err != nil {
		return participle.NextMatch
	}

	*node = parsed

	lex.Next()

	return nil
}

// UnOpSymbol represents binary operation symbols.
// ENUM(not)
type UnOpSymbol byte

func (node *UnOpSymbol) Parse(lex *lexer.PeekingLexer) error {
	token := lex.Peek()

	parsed, err := ParseUnOpSymbol(token.Value)
	if err != nil {
		return participle.NextMatch
	}

	*node = parsed

	lex.Next()

	return nil
}
