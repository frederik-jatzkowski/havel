package assembly

import (
	"encoding/binary"
	"fmt"
	"io"
)

type sLabel struct {
	label string
}

func (p *P) AddSLabel(label string) {
	p.StaticData = append(p.StaticData, &sLabel{label: label})
}

var _ S = &sLabel{}

func (s *sLabel) Bytes() int {
	return 8
}

func (s *sLabel) String(_ map[string]int) string {
	return fmt.Sprintf("  &%s", s.label)
}

func (s *sLabel) WriteTo(w io.Writer, labels map[string]int) (n int64, err error) {
	data := make([]byte, 8)

	binary.LittleEndian.PutUint64(data, uint64(labels[s.label]))

	length, err := w.Write(data)

	return int64(length), err
}
