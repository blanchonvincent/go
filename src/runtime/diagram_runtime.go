// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

func BuildRuntimeDiagram() Buffer {
	stopTheWorldGC("runtime diagram")

	var buff Buffer

	for _, p := range allp {
		buff.startEvent(DiagramP, p.id, p.status)

		m := p.m.ptr()
		if m != nil {
			buff.startEvent(DiagramM, m.id)

			if m.curg != nil {
				buff.startEvent(DiagramCurG, m.curg.goid, m.curg.stack.hi, m.curg.stack.lo)
			}
			if m.g0 != nil {
				buff.startEvent(DiagramG0, m.g0.goid, m.g0.stack.hi, m.g0.stack.lo)
			}
			if m.gsignal != nil {
				buff.startEvent(DiagramGSignal, m.gsignal.goid, m.gsignal.stack.hi, m.gsignal.stack.lo)
			}
		}

		num := p.runqtail - p.runqhead
		if p.runnext != 0 {
			num++
		}
		buff.startEvent(RunnableGoroutines, num)

		if p.runnext != 0 {
			gp := p.runnext.ptr()

			buff.writeInt64(gp.goid)
			buff.writeUintptr(gp.stack.hi)
			buff.writeUintptr(gp.stack.lo)
		}

		runqhead := p.runqhead
		runqtail := p.runqtail

		for runqhead != runqtail {
			runqtail--
			gp := p.runq[runqtail%uint32(len(p.runq))].ptr()

			buff.writeInt64(gp.goid)
			buff.writeUintptr(gp.stack.hi)
			buff.writeUintptr(gp.stack.lo)
		}

		buff.startEvent(WaitingSyncGoroutines, uint32(len(p.sudogcache)))
		for _, g := range p.sudogcache {
			gp := g.g

			buff.writeInt64(gp.goid)
			buff.writeUintptr(gp.stack.hi)
			buff.writeUintptr(gp.stack.lo)
		}
	}

	buff.startEvent(GlobalesGoroutines, sched.runqsize)

	currg := sched.runq.head.ptr()
	for {
		if currg == nil {
			break
		}
		buff.writeInt64(currg.goid)
		buff.writeUintptr(currg.stack.hi)
		buff.writeUintptr(currg.stack.lo)

		currg = currg.schedlink.ptr()
	}

	/* do something here */

	println("npidle", sched.npidle)
	currp := sched.pidle.ptr()
	for {
		if currp == nil {
			break
		}
		println("p idle", currp.id, "runnable", currp.runqtail-currp.runqhead)
		currp = currp.link.ptr()
	}

	println("nmidle", sched.nmidle)
	currm := sched.midle.ptr()
	for {
		if currm == nil {
			break
		}
		println("m idle", currm.id, "runnable", currm.curg)
		currm = currm.schedlink.ptr()
	}

	println("nmfreed", sched.nmfreed)
	println("nmspinning", sched.nmspinning)
	println("g dead", sched.gFree.n)

	startTheWorldGC()

	return buff
}
