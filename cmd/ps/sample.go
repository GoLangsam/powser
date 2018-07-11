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

// Ones are 1 1 1 1 1 ... = `1/(1-x)` with a simple pole at `x=1`.
func Ones() PS {
	return AdInfinitum(NewCoefficient(1, 1))
}

// Twos are 2 2 2 2 2 ... just for samples.
func Twos() PS {
	return AdInfinitum(NewCoefficient(2, 1))
}

// AdInfinitum repeates coefficient `c` ad infinitum
// and returns `c^i`.
func AdInfinitum(c Coefficient) PS {
	Z := New()
	go func(Z PS, c Coefficient) {
		defer Z.Close()
		for Z.Next() {
			Z.Send(c)
		}
	}(Z, c)
	return Z
}

func sample(n int) {

	switch n {

	// Plus Less
	case 1:
		fmt.Print("#", n, " Ones: ")
		Ones().Printn(20)
	case 2:
		fmt.Print("#", n, " Twos: ")
		Twos().Printn(20)
	case 3:
		fmt.Print("#", n, " Add : ")
		Ones().Plus(Twos()).Printn(20)
	case 4:
		fmt.Print("#", n, " 1+  : ")
		Ones().Plus().Printn(20)
	case 5:
		fmt.Print("#", n, " 1+2 : ")
		Ones().Plus(Twos()).Printn(20)
	case 6:
		fmt.Print("#", n, " 1+8 : ")
		Ones().Plus(Twos(), Ones(), Twos(), Ones(), Twos()).Printn(20)
	case 7:
		fmt.Print("#", n, " 1-  : ")
		Ones().Less().Printn(20)
	case 8:
		fmt.Print("#", n, " 1-2 : ")
		Ones().Less(Twos()).Printn(18)
	case 9:
		fmt.Print("#", n, " 1-9 : ")
		Ones().Less(Ones(), Twos(), Ones(), Twos(), Ones(), Twos()).Printn(18)

	// Plus Less - with short PS

	case 10:
		fmt.Print("#", n, " Quad: ")
		sqr().Printn(20)
	case 11:
		fmt.Print("#", n, " Poly: ")
		poly().Printn(20)
	case 12:
		fmt.Print("#", n, " Poly: ")
		poly().Printn(2)
	case 13:
		fmt.Print("#", n, " Add : ")
		Ones().Plus(poly()).Printn(17)
	case 14:
		fmt.Print("#", n, " 1+  : ")
		poly().Plus().Printn(20)
	case 15:
		fmt.Print("#", n, " 1+2 : ")
		poly().Plus(sqr()).Printn(17)
	case 16:
		fmt.Print("#", n, " 1+9 : ")
		poly().Plus(one(), lin(), sqr(), one(), lin(), sqr()).Printn(15)
	case 17:
		fmt.Print("#", n, " 1-  : ")
		sqr().Less().Printn(20)
	case 18:
		fmt.Print("#", n, " 1-2 : ")
		sqr().Less(lin()).Printn(18)
	case 19:
		fmt.Print("#", n, " 1-9 : ")
		poly().Less(one(), lin(), sqr(), one(), lin(), sqr()).Printn(18)

		// Add: (Into)-Methods? GetNextFrom, SendOneFrom, Append, GetWith ???

		// Multiply & Co
	case 20:
		fmt.Print("#", n, " CMul: ")
		Ones().CMul(aMinusOne()).Printn(18)
	case 21:
		fmt.Print("#", n, " XMul: ") // XMul = MonMul(1)
		Ones().XMul().Printn(20)
	case 22:
		fmt.Print("#", n, " Mul : ")
		Ones().Times(Ones()).Printn(20)
	case 23:
		fmt.Print("#", n, " 1*  : ")
		Ones().Times().Printn(20)
	case 24:
		fmt.Print("#", n, " 1*2 : ")
		Ones().Times(Twos()).Printn(16)
	case 25:
		fmt.Print("#", n, " 1*9 : ")
		Ones().Times(Ones(), Twos(), Ones(), Twos(), Ones(), Twos()).Printn(14)
	case 26:
		fmt.Print("#", n, " Subst 2in1: ")
		Ones().Subst(Twos()).Printn(6)
	case 27:
		fmt.Print("#", n, " MonSubst  : ")
		Ones().MonSubst(aMinusOne(), 3).Printn(14)
	case 28:
		fmt.Print("#", n, " Subst c...: ")
		Ones().Subst(AdInfinitum(aCoeff())).Printn(16)
	case 29:
		fmt.Print("#", n, " c*(-1)^4  : ")
		AdInfinitum(aCoeff()).MonSubst(aMinusOne(), 4).Printn(20)

		// Constructors:
	case 30:
		fmt.Print("#", n, " AdInfinitum(c)  : ")
		AdInfinitum(aCoeff()).Printn(20)
	case 31:
		fmt.Print("#", n, " Monomial(c, 4)  : ")
		Monomial(aCoeff(), 4).Printn(20)
	case 32:
		fmt.Print("#", n, " Binomial(c)     : ")
		Binomial(aCoeff()).Printn(16)
	case 33:
		fmt.Print("#", n, " Polynom(1,2,3,c): ")
		Polynom(aOne(), NewCoefficient(2, 1), NewCoefficient(3, 1), aCoeff()).Printn(16)

		// Cofficients:
	case 36:
		fmt.Print("#", n, " Shift twos by c : ")
		Twos().Shift(aCoeff()).Printn(11)

	case 37:
		fmt.Print("#", n, " MonMul(cn) * x^n: ")
		poly().MonMul(int(cn)).Printn(20)

		// Recip():
	case 38:
		fmt.Print("#", n, " 1 / (1+x)       : ")
		lin().Recip().Printn(14)

	case 39:
		fmt.Print("#", n, " 1 / Poly        : ")
		poly().Recip().Printn(9)

	case 40:
		fmt.Print("#", n, " 1 / Poly #2     : ")
		poly().Recip().Printn(2)

		// Analysis:
	case 41:
		fmt.Print("#", n, " Deriv: ")
		Ones().Deriv().Printn(11)
	case 42:
		fmt.Print("#", n, " Integ: ")
		Ones().Integ(aZero()).Printn(12)

	case 43:
		fmt.Print("#", n, " Deriv: ")
		sqr().Deriv().Printn(11)
	case 44:
		fmt.Print("#", n, " Integ: ")
		lin().Integ(aZero()).Printn(12)

		// Functions:
	case 45:
		fmt.Print("#", n, " Exp  : ")
		Ones().Exp().Printn(9)

	case 46:
		fmt.Print("#", n, " ATan : ")
		Ones().MonSubst(aMinusOne(), 2).Integ(aZero()).Printn(14)

	case 47:
		fmt.Print("#", n, " Exp  : ")
		lin().Exp().Printn(9)

	case 48:
		fmt.Print("#", n, " ATan : ")
		Ones().MonSubst(aMinusOne(), 2).Integ(aZero()).Printn(14)

	default:
		fmt.Println("No such sample #", n, " - max =", max)

	}

}

const max = 48 // # of samples

// ===========================================================================
