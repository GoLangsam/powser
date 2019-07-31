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

func sample(n int) {

	var c ps.Coefficient

	switch n {

	case 1:
		fmt.Print("#", n, " Poly: ")
		poly().Printn(10)

	case 2:
		fmt.Println("#", n, " Ones: ")
		O1, O2 := poly().Split()
		c, _ = O1.Get()
		fmt.Print("  O1 = ", c)
		c, _ = O1.Get()
		fmt.Print("  O1 = ", c)
		c, _ = O1.Get()
		fmt.Print("  O1 = ", c)
		fmt.Println()
		c, _ = O2.Get()
		fmt.Print("  O2 = ", c)
		c, _ = O2.Get()
		fmt.Print("  O2 = ", c)
		c, _ = O2.Get()
		fmt.Print("  O2 = ", c)

		fmt.Println()
		O1.Printn(10)
		fmt.Println(" O1 done")
		O2.Printn(10)
		fmt.Println(" O2 done")

	case 3:
		fmt.Println("#", n, " Ones: ")
		O1, O2 := poly().Split()
		c, _ = O1.Get()
		fmt.Print("  O1 = ", c)
		// O1.Drop()
		fmt.Println()
		c, _ = O2.Get()
		fmt.Print("  O2 = ", c)
		c, _ = O2.Get()
		fmt.Print("  O2 = ", c)
		c, _ = O2.Get()
		fmt.Print("  O2 = ", c)

		fmt.Println()
		O1.Printn(10)
		fmt.Println(" O1 done")
		O2.Printn(10)
		fmt.Println(" O2 done")

	case 4:
		fmt.Println("#", n, " Ones: ")
		O1, O2 := poly().Split()
		c, _ = O1.Get()
		fmt.Print("  O1 = ", c)
		O1.Drop()
		fmt.Println(" O1 done")
		c, _ = O2.Get()
		fmt.Print("  O2 = ", c)
		c, _ = O2.Get()
		fmt.Print("  O2 = ", c)
		c, _ = O2.Get()
		fmt.Print("  O2 = ", c)

		fmt.Println()
		O2.Printn(10)
		fmt.Println(" O2 done")

	default:
		fmt.Println("No such sample #", n, " - max =", max)

	}

}

const max = 4 // # of samples

// ===========================================================================
