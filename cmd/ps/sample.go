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

	switch n {

	// Series
	case 1:
		fmt.Print("#", n, " Ones = 1 / (1-x): ")
		ps.Ones().Printn(20)
	case 2:
		fmt.Print("#", n, " Harmonics       : ")
		ps.Harmonics().Printn(13)

	case 3:
		fmt.Print("#", n, " Factorials      : ")
		ps.Factorials().Printn(14)
	case 4:
		fmt.Print("#", n, " Fibonaccis      : ")
		ps.Fibonaccis().Printn(18)

	case 5:
		fmt.Print("#", n, " 1 / Fact = e^x  : ")
		ps.OneByFactorial().Printn(14)
	case 6:
		fmt.Print("#", n, " 1 / Fibonacci   : ")
		ps.OneByFibonacci().Printn(18)

	case 7:
		fmt.Print("#", n, " Sinus           : ")
		ps.Sin().Printn(14)
	case 8:
		fmt.Print("#", n, " Cos (-1)^i/(2i)!: ")
		ps.Cos().Printn(14)

	case 9:
		fmt.Print("#", n, " Sin^2 + Cos^2==1: ")
		S1, C1 := ps.Sincos()
		S2, C2 := ps.Sincos()
		S1.Times(S2).Plus(C1.Times(C2)).Printn(20)

	// Plus Less

	case 10:
		fmt.Print("#", n, " Ones: ")
		ps.Ones().Printn(20)
	case 11:
		fmt.Print("#", n, " Twos: ")
		ps.Twos().Printn(20)
	case 12:
		fmt.Print("#", n, " Twos: ")
		ps.Twos().Printn(2)
	case 13:
		fmt.Print("#", n, " Add : ")
		ps.Ones().Plus(ps.Twos()).Printn(20)
	case 14:
		fmt.Print("#", n, " 1+  : ")
		ps.Ones().Plus().Printn(20)
	case 15:
		fmt.Print("#", n, " 1+2 : ")
		ps.Ones().Plus(ps.Twos()).Printn(20)
	case 16:
		fmt.Print("#", n, " 1+8 : ")
		ps.Ones().Plus(ps.Twos(), ps.Ones(), ps.Twos(), ps.Ones(), ps.Twos()).Printn(20)
	case 17:
		fmt.Print("#", n, " 1-  : ")
		ps.Ones().Less().Printn(20)
	case 18:
		fmt.Print("#", n, " 1-2 : ")
		ps.Ones().Less(ps.Twos()).Printn(18)
	case 19:
		fmt.Print("#", n, " 1-9 : ")
		ps.Ones().Less(ps.Ones(), ps.Twos(), ps.Ones(), ps.Twos(), ps.Ones(), ps.Twos()).Printn(18)

	// Plus Less - with short PS

	case 20:
		fmt.Print("#", n, " Quad: ")
		sqr().Printn(20)
	case 21:
		fmt.Print("#", n, " Poly: ")
		poly().Printn(20)
	case 22:
		fmt.Print("#", n, " Poly: ")
		poly().Printn(2)
	case 23:
		fmt.Print("#", n, " Add : ")
		poly().Plus(ps.Ones()).Printn(17)
	case 24:
		fmt.Print("#", n, " 1+  : ")
		poly().Plus().Printn(20)
	case 25:
		fmt.Print("#", n, " 1+2 : ")
		poly().Plus(sqr()).Printn(17)
	case 26:
		fmt.Print("#", n, " 1+9 : ")
		poly().Plus(one(), lin(), sqr(), one(), lin(), sqr()).Printn(20)
	case 27:
		fmt.Print("#", n, " 1-  : ")
		sqr().Less().Printn(20)
	case 28:
		fmt.Print("#", n, " 1-2 : ")
		sqr().Less(lin()).Printn(20)
	case 29:
		fmt.Print("#", n, " 1-9 : ")
		poly().Less(one(), lin(), sqr(), one(), lin(), sqr()).Printn(18)

	// Multiply & Co

	case 30:
		fmt.Print("#", n, " CMul: ")
		ps.Ones().CMul(aMinusOne()).Printn(18)
	case 31:
		fmt.Print("#", n, " XMul: ") // XMul = MonMul(1)
		ps.Ones().XMul().Printn(20)
	case 32:
		fmt.Print("#", n, " Mul : ")
		ps.Ones().Times(ps.Ones()).Printn(20)
	case 33:
		fmt.Print("#", n, " 1*  : ")
		ps.Ones().Times().Printn(20)
	case 34:
		fmt.Print("#", n, " 1*2 : ")
		ps.Ones().Times(ps.Twos()).Printn(20)
	case 35:
		fmt.Print("#", n, " 1*9 : ")
		ps.Ones().Times(ps.Ones(), ps.Twos(), ps.Ones(), ps.Twos(), ps.Ones(), ps.Twos()).Printn(14)
	case 36:
		fmt.Print("#", n, " Subst 2in1: ")
		ps.Ones().Subst(ps.Twos()).Printn(14)
	case 37:
		fmt.Print("#", n, " MonSubst  : ")
		ps.Ones().MonSubst(aMinusOne(), 3).Printn(20)
	case 38:
		fmt.Print("#", n, " Subst c...: ")
		ps.Ones().Subst(ps.AdInfinitum(aCoeff())).Printn(16)
	case 39:
		fmt.Print("#", n, " c*(-1)^4  : ")
		ps.AdInfinitum(aCoeff()).MonSubst(aMinusOne(), 4).Printn(20)

		// Constructors:
	case 40:
		fmt.Print("#", n, " AdInfinitum(c)  : ")
		ps.AdInfinitum(aCoeff()).Printn(20)
	case 41:
		fmt.Print("#", n, " Monomial(c, 4)  : ")
		ps.Monomial(aCoeff(), 4).Printn(20)
	case 42:
		fmt.Print("#", n, " Binomial(c)     : ")
		ps.Binomial(aCoeff()).Printn(16)
	case 43:
		fmt.Print("#", n, " Polynom(1,2,3,c): ")
		ps.Polynom(aOne(), ps.NewCoefficient(2, 1), ps.NewCoefficient(3, 1), aCoeff()).Printn(16)

	case 44:
		fmt.Print("#", n, " Factorials      : ")
		ps.Factorials().Printn(10)
	case 45:
		fmt.Print("#", n, " Fibonaccis      : ")
		ps.Fibonaccis().Printn(14)

		// Cofficients:
	case 46:
		fmt.Print("#", n, " Shift twos by c : ")
		ps.Twos().Shift(aCoeff()).Printn(11)

	case 47:
		fmt.Print("#", n, " MonMul(cn) * x^n: ")
		poly().MonMul(int(cn)).Printn(20)

		// Recip():
	case 48:
		fmt.Print("#", n, " 1 / (1+x)       : ")
		lin().Recip().Printn(14)

	case 49:
		fmt.Print("#", n, " 1 / Poly        : ")
		poly().Recip().Printn(9)

	case 50:
		fmt.Print("#", n, " 1 / Poly #2     : ")
		poly().Recip().Printn(2)

		// Analysis:
	case 51:
		fmt.Print("#", n, " Deriv: ")
		ps.Ones().Deriv().Printn(18)
	case 52:
		fmt.Print("#", n, " Integ: ")
		ps.Ones().Integ(aZero()).Printn(18)

	case 53:
		fmt.Print("#", n, " Deriv: ")
		sqr().Deriv().Printn(14)
	case 54:
		fmt.Print("#", n, " Integ: ")
		lin().Integ(aZero()).Printn(14)

		// Functions:
	case 55:
		fmt.Print("#", n, " Exp  : ")
		ps.Ones().Exp().Printn(11)

	case 56:
		fmt.Print("#", n, " ATan : ")
		ps.Ones().MonSubst(aMinusOne(), 2).Integ(aZero()).Printn(18)

	case 57:
		fmt.Print("#", n, " Tan  : ")
		ps.Tan().Printn(14)

	case 58:
		fmt.Print("#", n, " Cot*x: ")
		ps.CotX().Printn(14)

	case 59:
		fmt.Print("#", n, " Sec  : ")
		ps.Sec().Printn(14)

	case 60:
		fmt.Print("#", n, " Csc*x: ")
		ps.CscX().Printn(14)

		// Add: (Into)-Methods? NextGetFrom, SendOneFrom, Append, GetWith ???

	default:
		fmt.Println("No such sample #", n, " - max =", max)

	}

}

const max = 60 // # of samples

// ===========================================================================
