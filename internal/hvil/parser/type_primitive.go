package parser

import (
	"strconv"

	"github.com/alecthomas/participle/v2/lexer"
)

type PrimitiveType struct {
	BitSize uint8 `parser:"@BitSize"`
	Tokens  []lexer.Token
}

func (t PrimitiveType) String() string {
	return strconv.FormatUint(uint64(t.BitSize), 10)
}

func (t PrimitiveType) Equals(other Type) bool {
	return false
}
