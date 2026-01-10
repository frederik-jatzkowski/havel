package bytecode

import (
	"fmt"
)

type R byte

func (r R) String() string {
	return fmt.Sprintf("r%d", r)
}
