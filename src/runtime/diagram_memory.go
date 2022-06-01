// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

func BuildMemoryDiagram() Buffer {
	stopTheWorldGC("memory diagram")

	var buff Buffer

	for _, p := range allp {
		buff.startEvent(DiagramP, p.id, p.status)

		m := p.m.ptr()
		if m != nil {
			buff.startEvent(DiagramM, m.id)

			if m.curg != nil {
				buff.startEvent(DiagramCurG, m.curg.goid, m.curg.stack.hi, m.curg.stack.lo)
			}
		}

		if p.mcache != nil {
			buff.startEvent(DiagramMCache)

			l := int64(0)
			for _, mspan := range p.mcache.alloc {
				if mspan != nil && mspan.allocCount > 0 {
					l++
				}
			}
			buff.writeInt64(l)

			for class, mspan := range p.mcache.alloc {
				if mspan != nil && mspan.allocCount > 0 {
					size := class_to_size[spanClass(class).sizeclass()]
					if mspan.spanclass.noscan() {
						buff.writeByte(byte(1))
					} else {
						buff.writeByte(byte(0))
					}
					buff.writeUint16(size)
					buff.writeUint16(mspan.allocCount)
					buff.writeUintptr(mspan.nelems)
				}
			}
		}
	}

	startTheWorldGC()

	return buff
}
