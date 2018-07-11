// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================
// Constructors - from coefficients.

// Conventions
// Upper-case for power series.
// Lower-case for coefficients.
// Input variables: U,V,...
// Output variables: ...,Y,Z

// Power-series constructors return channels on which power
// series flow. They start an encapsulated generator that
// puts the terms of the series on the channel.

// Monomial returns `c * x^n`.
func Monomial(c Coefficient, n int) PS {
	Z := New()
	go func(Z PS, c Coefficient, n int) {
		defer Z.Close()

		if IsZero(c) {
			return
		}

		for ; n > 0; n-- { // n-1 times aZero
			if !Z.Put(aZero()) {
				return
			}
		}
		Z.Put(c) // `c * x^n`

	}(Z, c, n)
	return Z
}

// Binomial returns `(1+x)^c`,
// a finite polynom iff `c` is a positive
// and an alternating infinite power series otherwise.
func Binomial(c Coefficient) PS {
	Z := New()
	go func(Z PS, c Coefficient) {
		defer Z.Close()

		i, iZ := 1, aOne() // `1`, `1/1`
		for !IsZero(c) {
			if !Z.Put(iZ) {
				return
			}
			iZ.Mul(iZ, aC().Mul(c, rat1byI(i))) // `iZ = iZ * c * 1/i`
			c.Sub(c, aOne())                    // `c = c-1`
			i++
		}
	}(Z, c)
	return Z
}

// Polynom converts coefficients, constant term `c` first,
// to a (finite) power series, the polynom in the coefficients.
func Polynom(a ...Coefficient) PS {
	Z := New()
	go func(Z PS, a ...Coefficient) {
		defer Z.Close()

		j := 0
		for j = len(a); j > 0; j-- {
			if !IsZero(a[j-1]) { // remove trailing zeros
				break
			}
		}

		for i := 0; i < j; i++ {
			if !Z.Put(a[i]) {
				return
			}
		}
	}(Z, a...)
	return Z
}

// ===========================================================================
