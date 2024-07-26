package errors

import (
	"io"

	"github.com/alecthomas/participle/v2/lexer"
)

type Collector struct {
	writer    io.Writer
	collected []HelpfulError
}

func NewCollector(writer io.Writer) *Collector {
	return &Collector{
		writer: writer,
	}
}

func (collector *Collector) Err(
	Pos lexer.Position,
	Name string,
	Short string,
) {
	err := HelpfulError{Pos: Pos, Name: Name, Short: Short}
	collector.collected = append(collector.collected, err)
	err.tryToWriteTo(collector.writer)
}

func (collector *Collector) HasErrors() bool {
	return len(collector.collected) > 0
}

func (collector *Collector) Errors() []HelpfulError {
	return collector.collected
}
