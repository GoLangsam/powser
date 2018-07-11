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

	var c Coefficient

	switch n {

	case 1:
		fmt.Print("#", n, " Poly: ")
		poly().Printn(10)

	case 2:
		fmt.Println("#", n, " Ones: ")
		OO := poly().Split()
		c, _ = OO[0].Get()
		fmt.Print("  OO[0] = ", c)
		c, _ = OO[0].Get()
		fmt.Print("  OO[0] = ", c)
		c, _ = OO[0].Get()
		fmt.Print("  OO[0] = ", c)
		fmt.Println()
		c, _ = OO[1].Get()
		fmt.Print("  OO[1] = ", c)
		c, _ = OO[1].Get()
		fmt.Print("  OO[1] = ", c)
		c, _ = OO[1].Get()
		fmt.Print("  OO[1] = ", c)

		fmt.Println()
		OO[0].Printn(10)
		fmt.Println(" OO[0] done")
		OO[1].Printn(10)
		fmt.Println(" OO[1] done")

	case 3:
		fmt.Println("#", n, " Ones: ")
		OO := poly().Split()
		c, _ = OO[0].Get()
		fmt.Print("  OO[0] = ", c)
		// OO[0].Drop()
		fmt.Println()
		c, _ = OO[1].Get()
		fmt.Print("  OO[1] = ", c)
		c, _ = OO[1].Get()
		fmt.Print("  OO[1] = ", c)
		c, _ = OO[1].Get()
		fmt.Print("  OO[1] = ", c)

		fmt.Println()
		OO[0].Printn(10)
		fmt.Println()
		OO[1].Printn(10)
		fmt.Println()

	case 4:
		fmt.Println("#", n, " Ones: ")
		OO := poly().Split()
		c, _ = OO[0].Get()
		fmt.Print("  OO[0] = ", c)
		OO[0].Drop()
		fmt.Println()
		c, _ = OO[1].Get()
		fmt.Print("  OO[1] = ", c)
		c, _ = OO[1].Get()
		fmt.Print("  OO[1] = ", c)
		c, _ = OO[1].Get()
		fmt.Print("  OO[1] = ", c)

		fmt.Println()
		OO[1].Printn(10)
		fmt.Println()

	default:
		fmt.Println("No such sample #", n, " - max =", max)

	}

}

const max = 4 // # of samples

// ===========================================================================
