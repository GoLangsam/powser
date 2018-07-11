// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package dchpair

import (
	. "github.com/GoLangsam/powser/dch"
)

type DchPair [2]*Dch

func NewPair() (pair DchPair) {
	pair[0] = New()
	pair[1] = New()
	return pair
}
