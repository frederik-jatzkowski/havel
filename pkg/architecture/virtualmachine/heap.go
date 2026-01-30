package virtualmachine

import "encoding/binary"

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
	// Slice once: This triggers a single bounds check for [offset : offset+2]
	// and provides a 2-byte sub-slice for the binary package.
	binary.LittleEndian.PutUint16(h.data[ptr.handle][ptr.offset:ptr.offset+2], value)
}

func (h *Heap) Store32(ptr FatPtr, value uint32) {
	binary.LittleEndian.PutUint32(h.data[ptr.handle][ptr.offset:ptr.offset+4], value)
}

func (h *Heap) Store64(ptr FatPtr, value uint64) {
	binary.LittleEndian.PutUint64(h.data[ptr.handle][ptr.offset:ptr.offset+8], value)
}

func (h *Heap) Load8(ptr FatPtr) uint8 {
	return h.data[ptr.handle][ptr.offset]
}

func (h *Heap) Load16(ptr FatPtr) uint16 {
	// The compiler will optimize this sub-slice and binary call
	// into a single MOV instruction on most modern architectures.
	return binary.LittleEndian.Uint16(h.data[ptr.handle][ptr.offset : ptr.offset+2])
}

func (h *Heap) Load32(ptr FatPtr) uint32 {
	return binary.LittleEndian.Uint32(h.data[ptr.handle][ptr.offset : ptr.offset+4])
}

func (h *Heap) Load64(ptr FatPtr) uint64 {
	return binary.LittleEndian.Uint64(h.data[ptr.handle][ptr.offset : ptr.offset+8])
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
		handle: uint32(value >> 32),
		offset: uint32(value),
	}
}

func (p FatPtr) ToUint64() uint64 {
	return (uint64(p.offset)) | uint64(p.handle)<<32
}
