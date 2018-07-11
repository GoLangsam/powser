// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big

type Rat struct {
	num, den int64 // numerator, denominator
}

func (u *Rat) Num() int64 { return u.num }
func (u *Rat) Den() int64 { return u.den }

func (u *Rat) Pr() {
	if u.den == 1 {
		print(u.num)
	} else {
		print(u.num, "/", u.den)
	}
	print(" ")
}

func (u *Rat) Eq(c *Rat) bool {
	return u.num == c.num && u.den == c.den
}

// Integer gcd; needed for rational arithmetic

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

// NewRatI creates a new Rat `i/1` from int `i`.
func NewRatI(i int) *Rat {
	return NewRat(int64(i), 1)
}

// NewRat1byI creates a new Rat `1/i` from int `i`.
func NewRat1byI(i int) *Rat {
	return NewRat(1, int64(i))
}

var Zero *Rat
var One *Rat
var Two *Rat
var MinusOne *Rat

// End mark and end test

var Finis *Rat

func init() {
	Zero = NewRat(0, 1)
	One = NewRat(1, 1)
	Two = NewRat(2, 1)
	MinusOne = Neg(One)
	Finis = NewRat(1, 0)
}

// End mark and end test
func (u *Rat) End() int64 {
	if u.den == 0 {
		return 1
	}
	return 0
}

// Operations on rationals

// Add sets z to the sum x+y and returns z.
func Add(x, y *Rat) *Rat {
	g := gcd(x.den, y.den)
	return NewRat(x.num*(y.den/g)+y.num*(x.den/g), x.den*(y.den/g))
}

// Mul sets z to the product x*y and returns z.
func Mul(x, y *Rat) *Rat {
	g1 := gcd(x.num, y.den)
	g2 := gcd(x.den, y.num)
	r := new(Rat)
	r.num = (x.num / g1) * (y.num / g2)
	r.den = (x.den / g2) * (y.den / g1)
	return r
}

// Neg sets z to -x and returns z.
func Neg(x *Rat) *Rat {
	return NewRat(-x.num, x.den)
}

// Sub sets z to the difference x-y and returns z.
func Sub(x, y *Rat) *Rat {
	return Add(x, Neg(y))
}

// Inv sets z to 1/x and returns z.
func Inv(x *Rat) *Rat { // invert a rat
	if x.num == 0 {
		panic("zero divide in inv")
	}
	return NewRat(x.den, x.num)
}
