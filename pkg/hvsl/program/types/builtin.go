package types

import (
	"context"
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
)

//go:generate go tool go-enum ./builtin.go

// Builtin represents the kind of the scope
// ENUM(bool, u8, u16, u32, u64)
type Builtin byte

func (node Builtin) Position() lexer.Position {
	return lexer.Position{}
}

func (node Builtin) Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

func (node Builtin) Wrap(err error) error {
	return err
}

func (node Builtin) Decl() *Decl {
	return &Decl{Name: node.String(), Typ: node}
}

func (node Builtin) ResolveNames(ctx context.Context) error {
	return nil
}

func (node Builtin) Equals(other Type) bool {
	builtin, ok := other.(Builtin)
	if !ok {
		return false
	}

	return builtin == node
}
