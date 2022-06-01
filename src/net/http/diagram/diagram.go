// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diagram

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func init() {
	http.HandleFunc("/debug/diagram/", Index)
}

// Index responds with the diagram buffer profile named by the request.
// For example, "/debug/diagram/runtime" serves the "runtime" diagram.
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="diagram"`)

	name := strings.TrimPrefix(r.URL.Path, "/debug/diagram/")
	if name == "runtime" {
		buffer := runtime.BuildRuntimeDiagram()
		w.Write(buffer.Buf)

		return
	}
	if name == "memory" {
		buffer := runtime.BuildMemoryDiagram()
		w.Write(buffer.Buf)

		return
	}
	serveError(w, http.StatusNotFound, "Unknown profile")
}

func serveError(w http.ResponseWriter, status int, txt string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Go-Diagram", "1")
	w.Header().Del("Content-Disposition")
	w.WriteHeader(status)
	fmt.Fprintln(w, txt)
}
