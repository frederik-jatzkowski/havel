package virtualmachine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFatPtr_layout(t *testing.T) {
	ptr := NewFatPtr(1, 1)

	reconstructed := NewFatPtrFromUint64(ptr.ToUint64() + 1)
	assert.Equal(t, uint32(1), reconstructed.handle)
	assert.Equal(t, uint32(2), reconstructed.offset)
}
