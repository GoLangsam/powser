// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================
// Consumers: Eval & Print

// Conventions
// Upper-case for power series.
// Lower-case for coefficients.
// Input variables: U,V,...
// Output variables: ...,Y,Z

// Eval n terms of power series `U` at `x=c`.
func (U PS) Eval(c Coefficient, n int) Coefficient {
	if n == 0 {
		U.Drop()
		return aZero
	}

	var u Coefficient
	var ok bool
	if u, ok = U.Get(); !ok {
		U.Drop()
		return aZero
	}
	return u.Add(u, c.Mul(c, U.Eval(c, n-1)))
}

// Evaln evaluates PS at `x=c` to n terms in floating point.
func (U PS) Evaln(c Coefficient, n int) float64 {
	defer U.Drop()

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
	return val
}

// Printn prints n terms of a power series.
func (U PS) Printn(n int) {
	defer U.Drop()
	defer print(("\n"))

	var u Coefficient
	var ok bool
	for ; n > 0; n-- {
		if u, ok = U.Get(); !ok {
			return
		}
		print(u.String())
		print(" ")
	}
}

// Print one billion terms.
func (U PS) Print() {
	U.Printn(1000000000)
}

// ===========================================================================
