// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package main

import (
	"fmt"
)

// ===========================================================================

func sample(n int) {

	switch n {

	case 1:
		fmt.Print("#", n, " Poly: ")
		poly().Printn(10)

	case 2:
		fmt.Println("#", n, " 1*1: ")
		one().Times(one()).Printn(N)
	case 3:
		fmt.Println("#", n, " 1+x: ")
		one().Times(lin()).Printn(N)
	case 4:
		fmt.Println("#", n, " x+1: ")
		lin().Times(one()).Printn(N)
	case 5:
		fmt.Println("#", n, " (1+x)^2: ")
		lin().Times(lin()).Printn(N)
	case 6:
		fmt.Println("#", n, " 1*3: ")
		one().Times(sqr()).Printn(N)
	case 7:
		fmt.Println("#", n, " 3*1: ")
		sqr().Times(one()).Printn(N)
	case 8:
		fmt.Println("#", n, " 2*3: ")
		lin().Times(sqr()).Printn(N)
	case 9:
		fmt.Println("#", n, " 3*2: ")
		sqr().Times(lin()).Printn(N)
	case 10:
		fmt.Println("#", n, " 3*3: ")
		sqr().Times(sqr()).Printn(N)
	default:
		fmt.Println("No such sample #", n, " - max =", max)

	}

}

const max = 10 // # of samples

// ===========================================================================
