// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps_test

import (
	ps "github.com/GoLangsam/powser"
)

// ===========================================================================
// special rational coefficients

// aC is just a shorthand useful in math with Coefficients.
func aC() ps.Coefficient {
	return ps.NewCoefficient(0, 1)
}

// aZero returns a zero.
func aZero() ps.Coefficient {
	return ps.NewCoefficient(0, 1)
}

// aOne returns a one.
func aOne() ps.Coefficient {
	return ps.NewCoefficient(1, 1)
}

// aMinusOne returns a minus one `-1`.
func aMinusOne() ps.Coefficient {
	return ps.NewCoefficient(-1, 1)
}

// ratIby1 creates a new Rat `i/1` from int `i`.
func ratIby1(i int) ps.Coefficient {
	return ps.NewCoefficient(int64(i), 1)
}

// rat1byI creates a new Rat `1/i` from int `i`.
func rat1byI(i int) ps.Coefficient {
	return ps.NewCoefficient(1, int64(i))
}

// ===========================================================================
