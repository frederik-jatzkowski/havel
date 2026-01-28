package virtualmachine

type Heap struct {
	data [][]byte
	free []uint32
}

func (h *Heap) Alloc(size uint64) FatPtr {
	if len(h.free) > 0 {
		handle := h.free[len(h.free)-1]
		h.free = h.free[:len(h.free)-1]

		h.data[handle] = make([]byte, size)

		return FatPtr{handle: handle}
	}

	h.data = append(h.data, make([]byte, size))

	return FatPtr{handle: uint32(len(h.data) - 1)}
}

func (h *Heap) Free(ptr FatPtr) {
	if h.data[ptr.handle] == nil {
		_ = h.data[ptr.handle][0] // intentional segfault
	}

	h.free = append(h.free, ptr.handle)
	h.data[ptr.handle] = nil
}

func (h *Heap) Store8(ptr FatPtr, value uint8) {
	h.data[ptr.handle][ptr.offset] = value
}

func (h *Heap) Store16(ptr FatPtr, value uint16) {
	h.data[ptr.handle][ptr.offset] = byte(value)
	h.data[ptr.handle][ptr.offset+1] = byte(value >> 8)
}

func (h *Heap) Store32(ptr FatPtr, value uint32) {
	h.data[ptr.handle][ptr.offset] = byte(value)
	h.data[ptr.handle][ptr.offset+1] = byte(value >> 8)
	h.data[ptr.handle][ptr.offset+2] = byte(value >> 16)
	h.data[ptr.handle][ptr.offset+3] = byte(value >> 24)
}

func (h *Heap) Store64(ptr FatPtr, value uint64) {
	h.data[ptr.handle][ptr.offset] = byte(value)
	h.data[ptr.handle][ptr.offset+1] = byte(value >> 8)
	h.data[ptr.handle][ptr.offset+2] = byte(value >> 16)
	h.data[ptr.handle][ptr.offset+3] = byte(value >> 24)
	h.data[ptr.handle][ptr.offset+1] = byte(value >> 32)
	h.data[ptr.handle][ptr.offset+2] = byte(value >> 40)
	h.data[ptr.handle][ptr.offset+3] = byte(value >> 48)
	h.data[ptr.handle][ptr.offset+1] = byte(value >> 56)
}

func (h *Heap) Load8(ptr FatPtr) uint8 {
	return h.data[ptr.handle][ptr.offset]
}

func (h *Heap) Load16(ptr FatPtr) uint16 {
	return uint16(h.data[ptr.handle][ptr.offset]) +
		(uint16(h.data[ptr.handle][ptr.offset+1]) << 8)
}

func (h *Heap) Load32(ptr FatPtr) uint32 {
	return uint32(h.data[ptr.handle][ptr.offset]) +
		(uint32(h.data[ptr.handle][ptr.offset+1]) << 8) +
		(uint32(h.data[ptr.handle][ptr.offset+1]) << 16) +
		(uint32(h.data[ptr.handle][ptr.offset+1]) << 24)
}

func (h *Heap) Load64(ptr FatPtr) uint64 {
	return uint64(h.data[ptr.handle][ptr.offset]) +
		(uint64(h.data[ptr.handle][ptr.offset+1]) << 8) +
		(uint64(h.data[ptr.handle][ptr.offset+1]) << 16) +
		(uint64(h.data[ptr.handle][ptr.offset+1]) << 24) +
		(uint64(h.data[ptr.handle][ptr.offset+1]) << 32) +
		(uint64(h.data[ptr.handle][ptr.offset+1]) << 40) +
		(uint64(h.data[ptr.handle][ptr.offset+1]) << 48) +
		(uint64(h.data[ptr.handle][ptr.offset+1]) << 56)
}

type FatPtr struct {
	handle uint32
	offset uint32
}

func NewFatPtr(handle uint32, offset uint32) FatPtr {
	return FatPtr{handle: handle, offset: offset}
}

func NewFatPtrFromUint64(value uint64) FatPtr {
	return FatPtr{
		handle: uint32(value),
		offset: uint32(value >> 32),
	}
}

func (p FatPtr) ToUint64() uint64 {
	return (uint64(p.offset) << 32) | uint64(p.handle)
}
