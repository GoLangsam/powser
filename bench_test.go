// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package ps_test

import (
	"testing"

	"github.com/GoLangsam/powser"
)

func BenchmarkPS_Times(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ps.Ones().Times(ps.Ones()).Printn(1000)
	}
}
