// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

import (
	"math/big"

	"github.com/GoLangsam/powser/dch"
	// "github.com/GoLangsam/powser/rat"
	// "github.com/GoLangsam/powser/rat.dch"
)

// ===========================================================================

// Coefficient of a power series - a rational number in this case.
//
//  Note: Coefficients just need to provide:
//  `Equal(a, b)` to discriminate iff `a` is equal to `b`.
//  `IsZero(a)` to discriminate iff coefficient `a` is equal to zero.
//  `Add(a, b)` as commutative and associative addition with `aZero` as neutral element,
//  `Sub(a, b)` as subtraction, and `Neg(a)` as convenience for `-a` (so `Add(a, Neg(a)) == aZero`),
//  `Mul(a, b)` as commutative and associative multiplication with `aOne` as neutral element,
//  `Inv(a)` as the inverse of multiplication `1/a` (for `a` not == aZero).
//
//  Note: `Inv(a)` is for `U.Recip()` only; remove for coefficients with no inverse.
type Coefficient = *big.Rat

// NewCoefficient returns a new coefficient: the rational `a/b`.
func NewCoefficient(a, b int64) Coefficient {
	return big.NewRat(a, b)
}

// IsZero discriminates iff coefficient `c` is equal to zero.
func IsZero(c Coefficient) bool {
	// return c.Num() == 0
	return c.Num().Cmp(big.NewInt(0)) == 0
}

// Equal discriminates iff coefficient `a` is equal to `b`.
func Equal(a, b Coefficient) bool {
	// return a.Num() == b.Num() && a.Denom() == b.Denom()
	return a.Cmp(b) == 0
}

// PS represents a power series as a demand channel
// of it's coefficients.
type PS struct {
	*dch.Dch
}

// New returns a fresh power series.
func New() PS {
	return PS{dch.DchMakeChan()}
}

// ===========================================================================
