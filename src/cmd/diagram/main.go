// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/blushft/go-diagrams/diagram"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"cmd/internal/objabi"
)

const usageMessage = "" +
	`Usage of 'go tool diagram':
Export a .dot file diagram for the runtime:
	go tool diagram -runtime

Export a .dot file diagram for the memory:
	go tool diagram -memory
`

func usage() {
	fmt.Fprint(os.Stderr, usageMessage)
	fmt.Fprintln(os.Stderr, "\nFlags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\n  Only one of -runtime or -memory may be set.")
	os.Exit(2)
}

var (
	input  = flag.String("i", "", "file for input")
	output  = flag.String("o", "", "file for output")
	runtimeOut = flag.Bool("runtime", false, "generate runtime representation")
	memoryOut = flag.Bool("memory", false, "generate memory representation")
)

func main() {
	objabi.AddVersionFlag()
	flag.Usage = usage
	flag.Parse()

	// Usage information when no arguments.
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.Usage()
	}

	err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, `For usage information, run "go tool diagram -help"`)
		os.Exit(2)
	}

	if *runtimeOut {
		d, err := diagram.New(diagram.Filename(*output), diagram.Label("Runtime"), diagram.Direction("LR"))
		if err != nil {
			log.Fatal(err)
		}

		pps := diagram.NewGroup("pps").Label("Ps").BackgroundColor("#efefefff")
		d.Group(pps)

		runtm := diagram.NewGroup("runtime").Label("Runtime").BackgroundColor("#efefefff")
		d.Group(runtm)

		var groupP *diagram.Group
		var nodeP, nodeM, nodeG, nodeG0, nodeGsignal *diagram.Node

		link := func() {
			defer func() {
				nodeP = nil
				nodeM = nil
				nodeG = nil
				nodeG0 = nil
				nodeGsignal = nil
			}()

			if groupP == nil {
				return
			}
			pps.Group(groupP)
			if nodeP == nil {
				return
			}
			groupP.Add(nodeP)
			if nodeM != nil {
				groupP.Add(nodeM)
				d.Connect(nodeP, nodeM, diagram.Forward())
			}
			if nodeG != nil {
				groupP.Add(nodeG)
				d.Connect(nodeM, nodeG, diagram.Forward())
			}
			if nodeG0 != nil {
				groupP.Add(nodeG0)
				d.Connect(nodeM, nodeG0, diagram.Forward())
			}
			if nodeGsignal != nil {
				groupP.Add(nodeGsignal)
				d.Connect(nodeM, nodeGsignal, diagram.Forward())
			}
		}

		for {
			e := b.NextEvent()
			if 0 == e {
				break
			}

			switch e {
			case runtime.DiagramP:
				link()

				id := fmt.Sprintf("P%d", b.ReadInt32())

				// status
				_ = b.ReadUint32()

				groupP = diagram.NewGroup(id).Label(id)
				nodeP = node.Runtime.P(diagram.NodeLabel(fmt.Sprintf("%s", id)))

			case runtime.DiagramM:
				id := b.ReadInt64()
				nodeM = node.Runtime.M(diagram.NodeLabel(fmt.Sprintf("M%d", id)))

			case runtime.DiagramCurG:
				id := b.ReadInt64()
				hi := b.ReadUintptr()
				lo := b.ReadUintptr()

				nodeG = node.Runtime.G(diagram.NodeLabel(fmt.Sprintf("G%d\nstack %dKB", id, (hi-lo)/(1<<10))))

			case runtime.DiagramG0:
				_ = b.ReadInt64()
				hi := b.ReadUintptr()
				lo := b.ReadUintptr()

				nodeG0 = node.Runtime.G(diagram.NodeLabel(fmt.Sprintf("g0\nstack %dKB", (hi-lo)/(1<<10))))

			case runtime.DiagramGSignal:
				_ = b.ReadInt64()
				hi := b.ReadUintptr()
				lo := b.ReadUintptr()

				nodeGsignal = node.Runtime.G(diagram.NodeLabel(fmt.Sprintf("g Signal\nstack %dKB", (hi-lo)/(1<<10))))

			case runtime.GlobalesGoroutines:
				ng := b.ReadInt32()
				if 0 == ng {
					continue
				}

				global := diagram.NewGroup("global").
					Label("Global goroutines").
					BackgroundColor("#e6b8af8a")

				maxDisplay := int32(3)
				var previous *diagram.Node
				for i := int32(0); i < ng; i++ {
					id := b.ReadInt64()
					hi := b.ReadUintptr()
					lo := b.ReadUintptr()

					g := node.Runtime.G(diagram.NodeLabel(fmt.Sprintf("G%d\nstack %dKB", id, (hi-lo)/(1<<10))))
					global.Add(g)

					if previous != nil {
						global.Connect(previous, g, diagram.Forward())
					}
					previous = g

					if i == (maxDisplay - 1) {
						num := ng - i - 1
						if num > 0 {
							more := node.Runtime.G(diagram.NodeLabel(fmt.Sprintf("%d more...", num)))
							global.Connect(previous, more, diagram.Forward())
						}

						break
					}
				}
				runtm.Group(global)

			case runtime.RunnableGoroutines:
				ng := b.ReadUint32()
				if 0 == ng {
					continue
				}

				runnable := diagram.NewGroup(fmt.Sprintf("runnable_%d", groupP)).
					Label("Runnable goroutines").
					BackgroundColor("#e6b8af8a")

				maxDisplay := uint32(3)
				var previous *diagram.Node
				for i := uint32(0); i < ng; i++ {
					id := b.ReadInt64()
					hi := b.ReadUintptr()
					lo := b.ReadUintptr()

					g := node.Runtime.G(diagram.NodeLabel(fmt.Sprintf("G%d\nstack %dKB", id, (hi-lo)/(1<<10))))
					runnable.Add(g)

					if previous != nil {
						runnable.Connect(previous, g, diagram.Forward())
					} else {
						groupP.Connect(nodeP, g, diagram.Forward())
					}
					previous = g

					if i == (maxDisplay - 1) {
						num := ng - i - 1
						if num > 0 {
							more := node.Runtime.G(diagram.NodeLabel(fmt.Sprintf("%d more...", num)))
							runnable.Connect(previous, more, diagram.Forward())
						}

						break
					}
				}
				groupP.Group(runnable)

			case runtime.WaitingSyncGoroutines:
				ng := b.ReadUint32()
				if 0 == ng {
					continue
				}

				waiting := diagram.NewGroup(fmt.Sprintf("waiting_sync_%d", groupP)).
					Label("Waiting goroutines for sync").
					BackgroundColor("#e6b8af8a")

				maxDisplay := uint32(3)
				var previous *diagram.Node
				for i := uint32(0); i < ng; i++ {
					id := b.ReadInt64()
					hi := b.ReadUintptr()
					lo := b.ReadUintptr()

					g := node.Runtime.G(diagram.NodeLabel(fmt.Sprintf("G%d\nstack %dKB", id, (hi-lo)/(1<<10))))
					waiting.Add(g)

					if previous != nil {
						waiting.Connect(previous, g, diagram.Forward())
					} else {
						groupP.Connect(nodeP, g, diagram.Forward())
					}
					previous = g

					if i == (maxDisplay - 1) {
						num := ng - i - 1
						if num > 0 {
							more := node.Runtime.G(diagram.NodeLabel(fmt.Sprintf("%d more...", num)))
							waiting.Connect(previous, more, diagram.Forward())
						}

						break
					}
				}
				groupP.Group(waiting)
			}
		}

		link()

		if err := d.Render(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	if *memoryOut {

	}
}

// parseFlags sets the profile and counterStmt globals and performs validations.
func parseFlags() error {
	if !*runtimeOut {
		if !*memoryOut {
			return fmt.Errorf("one of -runtime or -memory must be set")
		}
	}

	if *input == "" {
		return fmt.Errorf("you must specify the input with -i")
	}

	if *output == "" {
		return fmt.Errorf("you must specify the output with -o")
	}

	return nil
}

