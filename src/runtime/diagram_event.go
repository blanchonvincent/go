// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const (
	DiagramP int16 = iota + 1
	// id,status

	RunnableGoroutines
	// len,id,stack,id,stack,...

	GlobalesGoroutines
	// len,id,stack,id,stack,...

	WaitingSyncGoroutines
	// len,id,stack,id,stack,...

	DiagramM
	// id
	// go -> DiagramG
	// gsignals -> DiagramG

	DiagramMCache
	//

	DiagramCurG
	DiagramG0
	DiagramGSignal
	DiagramG
	// id,stack
)

type Buffer struct {
	readOff  uint32
	writeOff uint32
	Buf      []byte
	Events   []uint32
}

func (b *Buffer) NextEvent() int16 {
	if 0 == len(b.Events) {
		return 0
	}

	var event uint32
	event, b.Events = b.Events[0], b.Events[1:]
	b.readOff = event

	return b.ReadInt16()
}

func (b *Buffer) startEvent(v int16, data ...interface{}) {
	b.Events = append(b.Events, b.writeOff)
	b.writeInt16(v)

	for _, d := range data {
		switch d.(type) {
		case int32:
			b.writeInt32(d.(int32))
		case uint32:
			b.writeUint32(d.(uint32))
		case int64:
			b.writeInt64(d.(int64))
		case uintptr:
			b.writeUintptr(d.(uintptr))
		}
	}
}

func (b *Buffer) writeByte(v byte) {
	b.Buf = append(b.Buf, v)
	b.writeOff += 1
}

func (b *Buffer) ReadByte() byte {
	v := b.Buf[b.readOff]
	b.readOff += 1

	return v
}

func (b *Buffer) writeInt16(v int16) {
	b.Buf = append(b.Buf, byte(v), byte(v>>8))
	b.writeOff += 2
}

func (b *Buffer) ReadInt16() int16 {
	v := int16(b.Buf[b.readOff]) | int16(b.Buf[b.readOff+1])<<8
	b.readOff += 2

	return v
}

func (b *Buffer) writeUint16(v uint16) {
	b.Buf = append(b.Buf, byte(v), byte(v>>8))
	b.writeOff += 2
}

func (b *Buffer) ReadUint16() uint16 {
	v := uint16(b.Buf[b.readOff]) | uint16(b.Buf[b.readOff+1])<<8
	b.readOff += 2

	return v
}

func (b *Buffer) writeInt32(v int32) {
	b.Buf = append(b.Buf, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	b.writeOff += 4
}

func (b *Buffer) ReadInt32() int32 {
	v := int32(b.Buf[b.readOff]) | int32(b.Buf[b.readOff+1])<<8 | int32(b.Buf[b.readOff+2])<<16 | int32(b.Buf[b.readOff+3])<<24
	b.readOff += 4

	return v
}

func (b *Buffer) writeUint32(v uint32) {
	b.Buf = append(b.Buf, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	b.writeOff += 4
}

func (b *Buffer) ReadUint32() uint32 {
	v := uint32(b.Buf[b.readOff]) | uint32(b.Buf[b.readOff+1])<<8 | uint32(b.Buf[b.readOff+2])<<16 | uint32(b.Buf[b.readOff+3])<<24
	b.readOff += 4

	return v

}

func (b *Buffer) writeInt64(v int64) {
	b.Buf = append(b.Buf,
		byte(v), byte(v>>8), byte(v>>16), byte(v>>24),
		byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56),
	)
	b.writeOff += 8
}

func (b *Buffer) ReadInt64() int64 {
	v := int64(b.Buf[b.readOff]) | int64(b.Buf[b.readOff+1])<<8 | int64(b.Buf[b.readOff+2])<<16 | int64(b.Buf[b.readOff+3])<<24 | int64(b.Buf[b.readOff+4])<<32 | int64(b.Buf[b.readOff+5])<<40 | int64(b.Buf[b.readOff+6])<<48 | int64(b.Buf[b.readOff+7])<<56
	b.readOff += 8

	return v
}

func (b *Buffer) writeUint64(v uint64) {
	b.Buf = append(b.Buf,
		byte(v), byte(v>>8), byte(v>>16), byte(v>>24),
		byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56),
	)
	b.writeOff += 8
}

func (b *Buffer) ReadUint64() uint64 {
	v := uint64(b.Buf[b.readOff]) | uint64(b.Buf[b.readOff+1])<<8 | uint64(b.Buf[b.readOff+2])<<16 | uint64(b.Buf[b.readOff+3])<<24 | uint64(b.Buf[b.readOff+4])<<32 | uint64(b.Buf[b.readOff+5])<<40 | uint64(b.Buf[b.readOff+6])<<48 | uint64(b.Buf[b.readOff+7])<<56
	b.readOff += 8

	return v
}

func (b *Buffer) writeUintptr(v uintptr) {
	b.Buf = append(b.Buf,
		byte(v), byte(v>>8), byte(v>>16), byte(v>>24),
		byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56),
	)
	b.writeOff += 8
}

func (b *Buffer) ReadUintptr() uintptr {
	v := uintptr(b.Buf[b.readOff]) | uintptr(b.Buf[b.readOff+1])<<8 | uintptr(b.Buf[b.readOff+2])<<16 | uintptr(b.Buf[b.readOff+3])<<24 | uintptr(b.Buf[b.readOff+4])<<32 | uintptr(b.Buf[b.readOff+5])<<40 | uintptr(b.Buf[b.readOff+6])<<48 | uintptr(b.Buf[b.readOff+7])<<56
	b.readOff += 8

	return v
}
