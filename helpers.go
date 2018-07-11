// Copyright 2009 The Go Authors. All rights reserved.
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
	done := false
	for ; !done && n > 0; n-- {
		u := U.Get()
		if atEnd(u) {
			done = true
		} else {
			print(u.String())
		}
	}
	print(("\n"))
}

// Print one billion terms.
func (U PS) Print() {
	U.Printn(1000000000)
}

// ===========================================================================
// Helpers

// GetVal2 each first value received from the two given power series.
func GetVal2(U, V PS) (u, v Coefficient) {
	pair := getValS(U, V)
	return pair[0], pair[1]
}

// Split returns a pair of power series identical to a given power series
func (U PS) Split() PS2 {
	UU := NewPS2()
	go UU.Split(U)
	return UU
}

// Append the coefficient from `from` to `U`.
func (U PS) Append(from PS) {
	req, in := U.Into()
	for {
		<-req
		in <- from.Get()
	}
}

// Eval n terms of power series U at x=c
func (U PS) Eval(c Coefficient, n int) Coefficient {
	if n == 0 {
		return aZero
	}
	y := U.Get()
	if atEnd(y) {
		return aZero
	}
	return y.Add(y, c.Mul(c, U.Eval(c, n-1)))
}

// Evaln evaluates PS at `x=c` to n terms in floating point.
func (U PS) Evaln(c Coefficient, n int) float64 {
	xn := float64(1)
	x := float64(c.Num()) / float64(c.Denom())
	val := float64(0)
	for i := 0; i < n; i++ {
		u := U.Get()
		if atEnd(u) {
			break
		}
		val = val + x*float64(u.Num())/float64(u.Denom())
		xn = xn * x
	}
	return val
}

// ===========================================================================
