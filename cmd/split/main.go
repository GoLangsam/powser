// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	. "github.com/GoLangsam/powser"
)

var ( //flags
	n int
	x bool

	cn, cd int64
)

func init() {
	flag.IntVar(&n, "n", 0, "# of sample: 0 = all")
	flag.BoolVar(&x, "x", true, "use os.Exit(1) to see leaking goroutines, if any")
	flag.Int64Var(&cn, "cn", 1, "Coefficient: Numerator")
	flag.Int64Var(&cd, "cd", 1, "Coefficient: Denominator")
	flag.Parse()
}

// ===========================================================================
// special rational coefficients

var (
	aZero     Coefficient // 0
	aOne      Coefficient // 1
	aMinusOne Coefficient // -1
)

func init() {
	aZero = NewCoefficient(0, 1)      // 0
	aOne = NewCoefficient(1, 1)       // 1
	aMinusOne = NewCoefficient(-1, 1) // -1 - aMinusOne.Neg(aOne) raises `nil` exception ?!?
}

// ===========================================================================

func poly() PS {
	return Polynom(
		NewCoefficient(1, 1),
		NewCoefficient(1, 2),
		NewCoefficient(1, 3),
		NewCoefficient(1, 4),
		NewCoefficient(1, 5),
	)
}

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

	default:
		fmt.Println("No such sample #", n, " - max =", max)

	}

}

const max = 10 // # of samples

func main() {

	if n > 0 {
		sample(n)
	} else {
		for i := 1; i <= max; i++ {
			sample(i)
		}
	}

	if x {
		fmt.Println("about to leave ...")
		<-time.After(time.Millisecond)
		os.Exit(1) // to see leaking goroutines, if any
	}

}
