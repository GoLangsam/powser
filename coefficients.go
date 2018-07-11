// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================
// special rational coefficients

var (
	aZero     Coefficient
	aOne      Coefficient
	aTwo      Coefficient
	aMinusOne Coefficient
)

func init() {
	aZero = NewCoefficient(0, 1)
	aOne = NewCoefficient(1, 1)
	aTwo = NewCoefficient(2, 1)
	aMinusOne = NewCoefficient(-1, 1) // aMinusOne.Neg(aOne) raises `nil` exception ?!?
}

// ratI creates a new Rat `i/1` from int `i`.
func ratI(i int) Coefficient {
	return NewCoefficient(int64(i), 1)
}

// rat1byI creates a new Rat `1/i` from int `i`.
func rat1byI(i int) Coefficient {
	return NewCoefficient(1, int64(i))
}

// ===========================================================================
