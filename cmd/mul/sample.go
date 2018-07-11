// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package main

import (
	"fmt"

	. "github.com/GoLangsam/powser"
)

// ===========================================================================

func sample(n int) {

	switch n {

	case 1:
		fmt.Print("#", n, " Poly: ")
		poly().Printn(10)

	case 2:
		fmt.Println("#", n, " 1*1: ")
		Mul(one(), one()).Printn(N)
	case 3:
		fmt.Println("#", n, " 1+x: ")
		Mul(one(), lin()).Printn(N)
	case 4:
		fmt.Println("#", n, " x+1: ")
		Mul(lin(), one()).Printn(N)
	case 5:
		fmt.Println("#", n, " (1+x)^2: ")
		Mul(lin(), lin()).Printn(N)
	case 6:
		fmt.Println("#", n, " 1*3: ")
		Mul(one(), sqr()).Printn(N)
	case 7:
		fmt.Println("#", n, " 3*1: ")
		Mul(sqr(), one()).Printn(N)
	case 8:
		fmt.Println("#", n, " 2*3: ")
		Mul(lin(), sqr()).Printn(N)
	case 9:
		fmt.Println("#", n, " 3*2: ")
		Mul(sqr(), lin()).Printn(N)
	case 10:
		fmt.Println("#", n, " 3*3: ")
		Mul(sqr(), sqr()).Printn(N)
	default:
		fmt.Println("No such sample #", n, " - max =", max)

	}

}

const max = 10 // # of samples

// ===========================================================================
