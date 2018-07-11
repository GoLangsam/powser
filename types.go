// Copyright 2017 Andreas Pannewitz. All rights reserved.
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

// Eq discriminates iff x is equal to y.
func Eq(x, y Coefficient) bool {
	// return x.Num() == y.Num() && x.Denom() == y.Denom()
	return x.Cmp(y) == 0
}

// Integer represents the result of r.Nom() & r.Denum().
type Integer = *big.Int

// IsZero discriminates iff x is equal to zero.
func IsZero(x Integer) bool {
	// return x == 0
	return x.Cmp(big.NewInt(0)) == 0
}

// PS represents a power series as a demand channel
// of it's rational coefficients.
type PS struct {
	*dch.Dch
}

// New returns a fresh power series.
func New() PS {
	return PS{dch.New()}
}

// ===========================================================================
