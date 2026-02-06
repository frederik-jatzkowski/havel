package assembly

import (
	"encoding/binary"
	"fmt"
	"io"
)

type sLit struct {
	size  int
	value uint64
}

func (p *P) AddSLit(size int, value uint64) {
	p.StaticData = append(p.StaticData, &sLit{size, value})
}

var _ S = &sLit{}

func (s *sLit) Bytes() int {
	return s.size
}

func (s *sLit) String() string {
	return fmt.Sprintf("  %d (%d byte)", s.value, s.size)
}

func (s *sLit) WriteTo(w io.Writer) (n int64, err error) {
	data := make([]byte, s.size)

	switch s.size {
	case 1:
		data[0] = byte(s.value)
	case 2:
		binary.LittleEndian.PutUint16(data, uint16(s.value))
	case 4:
		binary.LittleEndian.PutUint32(data, uint32(s.value))
	case 8:
		binary.LittleEndian.PutUint64(data, s.value)
	}

	length, err := w.Write(data)

	return int64(length), err
}
