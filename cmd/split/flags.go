// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
)

// ===========================================================================

var ( //flags
	n int
	x bool

	cn, cd int64
)

func init() {
	flag.IntVar(&n, "n", 0, "# of sample: 0 = all")
	flag.BoolVar(&x, "x", true, "use os.Exit(1) to see leaking goroutines, if any")
	flag.Int64Var(&cn, "cn", 1, "Coefficient: Numerator")
	flag.Int64Var(&cd, "cd", 1, "Coefficient: Denominator")
	flag.Parse()
}

// ===========================================================================
