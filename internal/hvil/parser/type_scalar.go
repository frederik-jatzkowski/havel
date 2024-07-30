package parser

import (
	"strconv"

	"github.com/alecthomas/participle/v2/lexer"
)

type ScalarType struct {
	BitSize uint8 `parser:"@BitSize"`
	Tokens  []lexer.Token
}

func (t ScalarType) String() string {
	return strconv.FormatUint(uint64(t.BitSize), 10)
}

func (t ScalarType) Equals(other Type) bool {
	return t.String() == other.String()
}
