// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/GoLangsam/powser"
)

// ===========================================================================

// N may be used as "how-many-terms-to-print"
const N = 10

func aCoeff() ps.Coefficient {
	return ps.NewCoefficient(cn, cd)
}

func poly() ps.PS {
	return ps.Polynom(
		ps.NewCoefficient(1, 1),
		ps.NewCoefficient(1, 2),
		ps.NewCoefficient(1, 3),
		ps.NewCoefficient(1, 4),
		ps.NewCoefficient(1, 5),
	)
}

func one() ps.PS {
	return ps.Polynom(
		ps.NewCoefficient(1, 1),
	)
}

func lin() ps.PS {
	return ps.Polynom(
		ps.NewCoefficient(1, 1),
		ps.NewCoefficient(1, 1),
	)
}

func sqr() ps.PS {
	return ps.Polynom(
		ps.NewCoefficient(1, 1),
		ps.NewCoefficient(1, 1),
		ps.NewCoefficient(1, 1),
	)
}

// ===========================================================================
