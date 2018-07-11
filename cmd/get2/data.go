// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	. "github.com/GoLangsam/powser"
)

// ===========================================================================

const N = 10

func aCoeff() Coefficient {
	return NewCoefficient(cn, cd)
}

func poly() PS {
	return Polynom(
		NewCoefficient(1, 1),
		NewCoefficient(1, 2),
		NewCoefficient(1, 3),
		NewCoefficient(1, 4),
		NewCoefficient(1, 5),
	)
}

func one() PS {
	return Polynom(
		NewCoefficient(1, 1),
	)
}

func lin() PS {
	return Polynom(
		NewCoefficient(1, 1),
		NewCoefficient(1, 1),
	)
}

func sqr() PS {
	return Polynom(
		NewCoefficient(1, 1),
		NewCoefficient(1, 1),
		NewCoefficient(1, 1),
	)
}

// ===========================================================================
