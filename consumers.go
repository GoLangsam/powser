// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

import (
	"fmt"
)

// ===========================================================================
// Consumers: Eval & Print

// Conventions
// Upper-case for power series.
// Lower-case for coefficients.
// Input variables: U,V,...
// Output variables: ...,Y,Z

// EvalAt evaluates a power series at `x=c`
// for up to `n` terms.
// Note: n=1 denotes the first, the constant term.
func (U PS) EvalAt(c Coefficient, n int) Coefficient {
	u, ok := U.Get()
	switch {
	case ok && n == 1:
		return u
	case ok && n > 1:
		return aC().Add(u, aC().Mul(c, U.EvalAt(c, n-1))) // `u + c*U`
	default:
		U.Drop()
		return aZero()
	}
}

// EvalN evaluates a power series at `x=c`
// for up to `n` terms in floating point.
// Note: n=1 denotes the first, the constant term.
func (U PS) EvalN(c Coefficient, n int) float64 {
	ci := float64(1)
	fc, _ := c.Float64()
	val := float64(0)

	var u Coefficient
	var fu float64
	var ok bool
	for i := 0; i < n; i++ {
		if u, ok = U.Get(); !ok {
			break
		}
		fu, _ = u.Float64()
		val += fu * ci // val += `u(i) * c^i`
		ci = fc * ci   // `c^(i+1) = c * c^i`
	}
	U.Drop()
	return val
}

// Printn prints up to n terms of a power series.
func (U PS) Printn(n int) {
	defer fmt.Print("\n")

	var u Coefficient
	var ok bool
	for ; n > 0; n-- {
		if u, ok = U.Get(); !ok {
			break
		}
		fmt.Print(u.String())
		fmt.Print(" ")
	}
	U.Drop()
}

// Printer returns a copy of `U`,
// and concurrently prints up to n terms of it.
// Useful to inspect formulas as it can be chained.
func (U PS) Printer(n int) PS {
	UU := U.Split()

	go func(U PS, n int) {
		U.Printn(n)
	}(U, n)
	return UU[1]
}

// Print one billion terms. Use at Your own risk ;-)
func (U PS) Print() {
	U.Printn(1000000000)
}

// ===========================================================================
