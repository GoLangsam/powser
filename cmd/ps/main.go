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

	aCoeff Coefficient // flag
)

func init() {
	aZero = NewCoefficient(0, 1)      // 0
	aOne = NewCoefficient(1, 1)       // 1
	aMinusOne = NewCoefficient(-1, 1) // -1 - aMinusOne.Neg(aOne) raises `nil` exception ?!?

	aCoeff = NewCoefficient(cn, cd)
}

// ===========================================================================

const N = 10

func sample(n int) {

	switch n {

	// Add Plus Less Minus
	case 1:
		fmt.Print("#", n, " Ones: ")
		Ones().Printn(20)
	case 2:
		fmt.Print("#", n, " Twos: ")
		Twos().Printn(20)
	case 3:
		fmt.Print("#", n, " Add : ")
		Add(Ones(), Twos()).Printn(17)
	case 4:
		fmt.Print("#", n, " 1+  : ")
		Ones().Plus().Printn(20)
	case 5:
		fmt.Print("#", n, " 1+2 : ")
		Ones().Plus(Twos()).Printn(17)
	case 6:
		fmt.Print("#", n, " 1+9 : ")
		Ones().Plus(Ones(), Twos(), Ones(), Twos(), Ones(), Twos()).Printn(15)
	case 7:
		fmt.Print("#", n, " 1-  : ")
		Ones().Less().Printn(20)
	case 8:
		fmt.Print("#", n, " 1-2 : ")
		Ones().Minus(Twos()).Printn(18)
	case 9:
		fmt.Print("#", n, " 1-9 : ")
		Ones().Less(Ones(), Twos(), Ones(), Twos(), Ones(), Twos()).Printn(18)

		// TODO: Append, clone

		// Multiply & Co
		// TODO: Monmul
		// TODO: Recip
	case 10:
		fmt.Print("#", n, " CMul: ")
		Ones().Cmul(aMinusOne).Printn(18)
	case 11:
		fmt.Print("#", n, " XMul: ") // Xmul = MonMul(1)
		Ones().Xmul().Printn(20)
	case 12:
		fmt.Print("#", n, " Mul : ")
		Mul(Ones(), Ones()).Printn(5)
	case 13:
		fmt.Print("#", n, " 1*  : ")
		Ones().Times().Printn(20)
	case 14:
		fmt.Print("#", n, " 1*2 : ")
		Ones().Times(Twos()).Printn(6)
	case 15:
		fmt.Print("#", n, " 1*9 : ")
		Ones().Times(Ones(), Twos(), Ones(), Twos(), Ones(), Twos()).Printn(4)

	case 16:
		fmt.Print("#", n, " Subst 2in1: ")
		Ones().Subst(Twos()).Printn(6)
	case 17:
		fmt.Print("#", n, " MonSubst  : ")
		Ones().MonSubst(aMinusOne, 3).Printn(14)
	case 18:
		fmt.Print("#", n, " Subst c...: ")
		Ones().Subst(AdInfinitum(aCoeff)).Printn(4)
	case 19:
		fmt.Print("#", n, " MonSubst U(c * (-1)^4) : ")
		AdInfinitum(aCoeff).MonSubst(aMinusOne, 4).Printn(10)

		// Constructors:
	case 20:
		fmt.Print("#", n, " AdInfinitum(c)  : ")
		AdInfinitum(aCoeff).Printn(11)
	case 21:
		fmt.Print("#", n, " Monomial(c,10)  : ")
		Monomial(aCoeff, 10).Printn(11)
	case 22:
		fmt.Print("#", n, " Binomial(c)     : ")
		Binomial(aCoeff).Printn(11)
	case 23:
		fmt.Print("#", n, " Polynom(1,2,3,c): ")
		Polynom(aOne, NewCoefficient(2, 1), NewCoefficient(3, 1), aCoeff).Printn(11)

		// Cofficients:
	case 26:
		fmt.Print("#", n, " Shift twos by c : ")
		Twos().Shift(aCoeff).Printn(11)

		// Analysis:
	case 31:
		fmt.Print("#", n, " Deriv: ")
		Ones().Deriv().Printn(11)
	case 32:
		fmt.Print("#", n, " Integ: ")
		Ones().Integ(aZero).Printn(12)

		// Functions:
	case 35:
		fmt.Print("#", n, " Exp: ")
		Ones().Exp().Printn(6)

	case 36:
		fmt.Print("#", n, " ATan: ")
		Ones().MonSubst(aMinusOne, 2).Integ(aZero).Printn(14)

	default:
		fmt.Println("No such sample #", n, " - max =", max)

	}

}

const max = 36 // # of samples

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
