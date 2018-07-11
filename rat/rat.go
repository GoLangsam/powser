// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rat

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

// Make a rational from two ints and from one int

func I2toR(u, v int64) *Rat {
	g := gcd(u, v)
	r := new(Rat)
	if v > 0 {
		r.num = u / g
		r.den = v / g
	} else {
		r.num = -u / g
		r.den = -v / g
	}
	return r
}

func ItoR(u int64) *Rat {
	return I2toR(u, 1)
}

var Zero *Rat
var One *Rat
var MinusOne *Rat

// End mark and end test

var Finis *Rat

func init() {
	Zero = ItoR(0)
	One = ItoR(1)
	MinusOne = Neg(One)
	Finis = I2toR(1, 0)
}

// End mark and end test
func (u *Rat) End() int64 {
	if u.den == 0 {
		return 1
	}
	return 0
}

// Operations on rationals

func Add(u, v *Rat) *Rat {
	g := gcd(u.den, v.den)
	return I2toR(u.num*(v.den/g)+v.num*(u.den/g), u.den*(v.den/g))
}

func Mul(u, v *Rat) *Rat {
	g1 := gcd(u.num, v.den)
	g2 := gcd(u.den, v.num)
	r := new(Rat)
	r.num = (u.num / g1) * (v.num / g2)
	r.den = (u.den / g2) * (v.den / g1)
	return r
}

func Neg(u *Rat) *Rat {
	return I2toR(-u.num, u.den)
}

func Sub(u, v *Rat) *Rat {
	return Add(u, Neg(v))
}

func Inv(u *Rat) *Rat { // invert a rat
	if u.num == 0 {
		panic("zero divide in inv")
	}
	return I2toR(u.den, u.num)
}
