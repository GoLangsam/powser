// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package big implements arithmetic on rational numbers
// as a subset of the methods of "math/big":
//  Num, Denom, Set, String
//  Add, Mul, Sub, Neg.
// Nominator and Denominator are int64.
package big

// A Rat represents a quotient a/b of arbitrary precision.
// The zero value for a Rat represents the value 0.
type Rat struct {
	num, den int64 // numerator, denominator
}

// Num returns the numerator of x; it may be <= 0.
func (x *Rat) Num() int64 { return x.num }

// Denom returns the denominator of x; it is always > 0.
func (x *Rat) Denom() int64 { return x.den }

// gcd needed for rational arithmetic
func gcd(u, v int64) int64 {
	if u < 0 {
		return gcd(-u, v)
	}
	if u == 0 {
		return v
	}
	return gcd(v%u, u)
}

// NewRat creates a new Rat with numerator a and denominator b.
func NewRat(a, b int64) *Rat {
	g := gcd(a, b)
	r := new(Rat)
	if b > 0 {
		r.num = a / g
		r.den = b / g
	} else {
		r.num = -a / g
		r.den = -b / g
	}
	return r
}

// Set sets `z` to `x` (by making a copy of `x`) and returns `z`.
func (z *Rat) Set(x *Rat) *Rat {
	if z == nil {
		z = new(Rat)
	}
	if z != x {
		z.num, z.den = x.num, x.den
	}
	return z
}

// Operations on rationals

// Add sets z to the sum x+y and returns z.
func (z *Rat) Add(x, y *Rat) *Rat {
	g := gcd(x.den, y.den)

	z.Set(x)
	z.num, z.den = x.num*(y.den/g)+y.num*(x.den/g), x.den*(y.den/g)
	return z
}

// Mul sets z to the product x*y and returns z.
func (z *Rat) Mul(x, y *Rat) *Rat {
	g1 := gcd(x.num, y.den)
	g2 := gcd(x.den, y.num)

	z.Set(x)
	z.num = (x.num / g1) * (y.num / g2)
	z.den = (x.den / g2) * (y.den / g1)
	return z
}

// Neg sets z to -x and returns z.
func (z *Rat) Neg(x *Rat) *Rat {
	z.Set(x)
	z.num, z.den = -x.num, x.den
	return z
}

// Sub sets z to the difference x-y and returns z.
func (z *Rat) Sub(x, y *Rat) *Rat {
	z.Set(x)
	return z.Add(x, z.Neg(y))
}

// Inv sets z to 1/x and returns z.
func (z *Rat) Inv(x *Rat) *Rat { // invert a rat
	if x.num == 0 {
		panic("zero divide in inv")
	}
	z.Set(x)
	z.num, z.den = x.den, x.num
	return z
}
