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

// AdInfinitum repeates coefficient `c` ad infinitum
// and returns `c^i`.
func AdInfinitum(c Coefficient) PS {
	Z := NewPS()
	go func(Z PS, c Coefficient) {
		defer Z.Close()
		for Z.Next() {
			Z.Send(c)
		}
	}(Z, c)
	return Z
}

// the Monomial of the coefficient
// returns `c * x^n`.
func Monomial(c Coefficient, n int) PS {
	Z := NewPS()
	go func(Z PS, c Coefficient, n int) {
		defer Z.Close()

		if !isZero(c.Num()) {
			for ; n > 0; n-- {
				if !Z.Put(aZero) {
					return
				}
			}
			Z.Put(c)
		}
	}(Z, c, n)
	return Z
}

// the Binomial theorem is applied to the coefficient
// and returns `(1+x)^c`.
func Binomial(c Coefficient) PS {
	Z := NewPS()
	go func(Z PS, c Coefficient) {
		defer Z.Close()

		n := 1
		t := aOne
		for !isZero(c.Num()) {
			Z.Put(t)
			t.Mul(t.Mul(t, c), rat1byI(n))
			c.Sub(c, aOne)
			n++
		}
	}(Z, c)
	return Z
}

// Polynom converts coefficients, constant term `c` first,
// to a (finite) power series, the polynom in the coefficients.
func Polynom(a ...Coefficient) PS {
	Z := NewPS()
	go func(Z PS, a ...Coefficient) {
		defer Z.Close()

		var done bool
		j := 0
		for j = len(a); !done && j > 0; j-- {
			if !isZero(a[j-1].Num()) { // remove trailing zeros
				done = true
			}
		}

		for i := 0; i < j; i++ {
			Z.Put(a[i])
		}
	}(Z, a...)
	return Z
}

// ===========================================================================
