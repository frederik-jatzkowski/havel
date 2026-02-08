package assembly

import (
	"io"
)

type S interface {
	String(labels map[string]int) string
	WriteTo(w io.Writer, labels map[string]int) (n int64, err error)
	Bytes() int
}
