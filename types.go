// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

import (
	"math/big"

	// "github.com/GoLangsam/powser/big"
	// "github.com/GoLangsam/powser/dch.rat"
	"github.com/GoLangsam/powser/dch.big"
)

// ===========================================================================

// Coefficient of a power series - a rational number.
type Coefficient = *big.Rat

// NewCoefficient returns a new coefficient `a/b`.
func NewCoefficient(a, b int64) Coefficient {
	return big.NewRat(a, b)
}

// integer represents the result of r.Nom() & r.Denum()
type integer = *big.Int

func isZero(x integer) bool {
	// return x == 0
	return x.Cmp(big.NewInt(0)) == 0
}

func asFloat64(x integer) float64 {
	// return float64(x)
	return float64(x.Int64())
}

// PS represents a power series as a demand channel
// of it's rational coefficients.
type PS struct {
	*dch.Dch
}

// NewPS returns a fresh power series.
func NewPS() PS {
	return PS{dch.New()}
}

// PS2 represents a pair of power series.
type PS2 [2]PS

// NewPS2 returns an empty pair of new power series.
func NewPS2() PS2 {
	return PS2{NewPS(), NewPS()}
}

// ===========================================================================