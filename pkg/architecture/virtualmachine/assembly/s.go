package assembly

import (
	"fmt"
	"io"
)

type S interface {
	fmt.Stringer
	io.WriterTo
	Bytes() int
}
