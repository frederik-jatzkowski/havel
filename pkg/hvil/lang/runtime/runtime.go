package runtime

import (
	"github.com/alecthomas/participle/v2/lexer"
	"io"
)

type Call struct {
	Name string
	Pos  lexer.Position
}

type VirtualMachine struct {
	Stack          []byte
	StackPointer   int
	CallStack      []Call
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

func New(
	stackSize int,
	stdin io.Reader,
	stdout, stderr io.Writer,
) *VirtualMachine {
	return &VirtualMachine{
		Stack:        make([]byte, stackSize),
		StackPointer: 0,
		Stdin:        stdin,
		Stdout:       stdout,
		Stderr:       stderr,
	}
}
