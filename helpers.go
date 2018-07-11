// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================
// Helpers

// Conventions
// Upper-case for power series.
// Lower-case for coefficients.
// Input variables: U,V,...
// Output variables: ...,Y,Z

// Printn prints n terms of a power series.
func (U PS) Printn(n int) {
	defer print(("\n"))

	var u Coefficient
	var ok bool
	for ; n > 0; n-- {
		if u, ok = U.Get(); !ok {
			return
		}
		print(u.String())
	}
}

// Print one billion terms.
func (U PS) Print() {
	U.Printn(1000000000)
}

// ===========================================================================
// Helpers

// New returns a fresh power series.
func (U PS) New() PS {
	return NewPS()
}

// NewPair returns an empty pair of new power series.
func (U PS) NewPair() PS2 {
	return PS2{NewPS(), NewPS()}
}

// Append the coefficient from `Z` to `U`.
func (U PS) Append(Z PS) {

	var u Coefficient
	var ok bool
	for U.Next() {
		if u, ok = Z.Get(); !ok {
			return
		}
		U.Send(u)
	}
}

// Eval n terms of power series `U` at `x=c`.
func (U PS) Eval(c Coefficient, n int) Coefficient {
	if n == 0 {
		return aZero
	}

	var u Coefficient
	var ok bool
	if u, ok = U.Get(); !ok {
		return aZero
	}
	return u.Add(u, c.Mul(c, U.Eval(c, n-1)))
}

// Evaln evaluates PS at `x=c` to n terms in floating point.
func (U PS) Evaln(c Coefficient, n int) float64 {
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

// ===========================================================================
