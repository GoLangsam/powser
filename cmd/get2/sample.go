// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package main

import (
	"fmt"

	"github.com/GoLangsam/powser"
)

// ===========================================================================

func test2(U, V ps.PS, n int) {

	for i := 0; i < n; i++ {
		u, okU, v, okV := U.GetWith(V)
		fmt.Println("\tu & v - Test #", i, "\tu", u, okU, "\tv", v, okV)
	}
	U.Drop()
	V.Drop()
}

func sample(n int) {

	switch n {

	case 1:
		fmt.Println("#", n, " 1 - Poly: ")
		test2(one(), poly(), 10)

	case 2:
		fmt.Println("#", n, " 2 - Poly: ")
		test2(lin(), poly(), 10)

	case 3:
		fmt.Println("#", n, " 1 - Poly: only first: ")
		test2(one(), poly(), 1)

	case 4:
		fmt.Println("#", n, " 2 - Poly: only first: ")
		test2(lin(), poly(), 1)

	default:
		fmt.Println("No such sample #", n, " - max =", max)

	}

}

const max = 4 // # of samples

// ===========================================================================
