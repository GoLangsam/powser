// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

import (
	"github.com/GoLangsam/powser/big"
	"github.com/GoLangsam/powser/dch"
)

// ===========================================================================

// PS represents a power series as a demand channel
// of it's rational coefficients.
type PS struct {
	*dch.Dch
}

// NewPS returns a new power series.
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

// Conventions
// Upper-case for power series.
// Lower-case for rationals.
// Input variables: U,V,...
// Output variables: ...,Y,Z

// Printn prints n terms of a power series
func (U PS) Printn(n int) {
	done := false
	for ; !done && n > 0; n-- {
		u := U.Get()
		if u.End() != 0 {
			done = true
		} else {
			print(u.String())
		}
	}
	print(("\n"))
}

// Print one billion terms
func (U PS) Print() {
	U.Printn(1000000000)
}

// ===========================================================================
// Helpers

// GetValS returns a slice with each first value received from the given power series.
func GetValS(U, V PS) (u, v *big.Rat) {
	pair := PS2{U, V}.GetValS() // todo: its akward - improve
	return pair[0], pair[1]
}

// Repeat keeps sending `dat` into `U`
func (U PS) Repeat(dat *big.Rat) {
	for {
		U.Put(dat)
	}
}

// Split returns a pair of power series identical to a given power series
func (U PS) Split() PS2 {
	UU := NewPS2()
	go UU.Split(U)
	return UU
}

// Copy data from `from` into `into`
func (U PS) Copy(from PS) {
	req, in := U.Into()
	for {
		<-req
		in <- from.Get()
	}
}

// Eval n terms of power series U at x=c
func (U PS) Eval(c *big.Rat, n int) *big.Rat {
	if n == 0 {
		return big.Zero
	}
	y := U.Get()
	if y.End() != 0 {
		return big.Zero
	}
	return y.Add(y, c.Mul(c, U.Eval(c, n-1)))
}

// Evaln evaluates PS at `x=c` to n terms in floating point.
func (U PS) Evaln(c *big.Rat, n int) float64 {
	xn := float64(1)
	x := float64(c.Num()) / float64(c.Denom())
	val := float64(0)
	for i := 0; i < n; i++ {
		u := U.Get()
		if u.End() != 0 {
			break
		}
		val = val + x*float64(u.Num())/float64(u.Denom())
		xn = xn * x
	}
	return val
}

// ===========================================================================
// ???

// Binom ial theorem
// and returns `(1+x)^c`
func Binom(c *big.Rat) PS {
	Z := NewPS()
	go func(c *big.Rat, Z PS) {
		n := 1
		t := big.One
		for c.Num() != 0 {
			Z.Put(t)
			t.Mul(t.Mul(t, c), big.NewRat1byI(n))
			c.Sub(c, big.One)
			n++
		}
		Z.Put(big.Finis)
	}(c, Z)
	return Z
}

// ===========================================================================
// Arithmetic on power series: each spawns a goroutine

// Power-series constructors return channels on which power
// series flow.  They start an encapsulated generator that
// puts the terms of the series on the channel.

// Add two power series
func Add(U, V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		Z_req, Z_in := Z.Into()
		var u, v *big.Rat
		for {
			<-Z_req
			u, v = GetValS(U, V)
			switch u.End() + 2*v.End() {
			case 0:
				Z_in <- u.Add(u, u)
			case 1:
				Z_in <- v
				Z.Copy(V)
			case 2:
				Z_in <- u
				Z.Copy(U)
			case 3:
				Z_in <- big.Finis
			}
		}
	}(U, V, Z)
	return Z
}

func (U PS) Plus(V PS) PS {
	return Add(U, V)
}

// Cmul multiplies a power series by a constant
func (U PS) Cmul(c *big.Rat) PS {
	Z := NewPS()
	go func(c *big.Rat, U, Z PS) {
		Z_req, Z_in := Z.Into()
		done := false
		for !done {
			<-Z_req
			u := U.Get()
			if u.End() != 0 {
				done = true
			} else {
				Z_in <- u.Mul(c, u)
			}
		}
		Z_in <- big.Finis
	}(c, U, Z)
	return Z
}

// Minus subtracts `V` from `U`
// and returns `U + (-1)*V`
func (U PS) Minus(V PS) PS {
	return U.Plus(V.Cmul(big.MinusOne))
}

// Monmul multiplies a power series by the monomial "x^n"
// and returns `x^n * U`.
func (U PS) Monmul(n int) PS {
	Z := NewPS()
	go func(n int, U PS, Z PS) {
		for ; n > 0; n-- {
			Z.Put(big.Zero)
		}
		Z.Copy(U)
	}(n, U, Z)
	return Z
}

// Xmul multiplies a power series by x, (by the monomial "x^1")
// and returns `x * U`.
func (U PS) Xmul() PS {
	return U.Monmul(1)
}

// Rep repeates c
// and returns `c^i`.
func Rep(c *big.Rat) PS {
	Z := NewPS()
	go Z.Repeat(c)
	return Z
}

// Mon returns the Monomial `c * x^n`
func Mon(c *big.Rat, n int) PS {
	Z := NewPS()
	go func(c *big.Rat, n int, Z PS) {
		if c.Num() != 0 {
			for ; n > 0; n = n - 1 {
				Z.Put(big.Zero)
			}
			Z.Put(c)
		}
		Z.Put(big.Finis)
	}(c, n, Z)
	return Z
}

// Shift
func (U PS) Shift(c *big.Rat) PS {
	Z := NewPS()
	go func(c *big.Rat, U, Z PS) {
		Z.Put(c)
		Z.Copy(U)
	}(c, U, Z)
	return Z
}

// simple pole at 1: 1/(1-x) = 1 1 1 1 1 ...

// Poly converts coefficients, constant term first
// to a (finite) power series, a polynom.
func Poly(a ...*big.Rat) PS {
	Z := NewPS()
	go func(Z PS, a ...*big.Rat) {
		var done bool
		j := 0
		for j = len(a); !done && j > 0; j-- {
			if a[j-1].Num() != 0 {
				done = true
			}
		}

		for i := 0; i < j; i++ {
			Z.Put(a[i])
		}

		Z.Put(big.Finis)
	}(Z, a...)
	return Z
}

// Mul multiplies. The algorithm is
//	let U = `u + x*UU`
//	let V = `v + x*VV`
//	then UV = `u*v + x*(u*VV+v*UU) + x*x*UU*VV`
func Mul(U, V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		Z_req, Z_in := Z.Into()
		<-Z_req
		u, v := GetValS(U, V)
		if u.End() != 0 || v.End() != 0 {
			Z_in <- big.Finis
		} else {
			var prod *big.Rat
			prod.Mul(u, v)
			Z_in <- prod
			UU := U.Split()
			VV := V.Split()
			W := Add(VV[0].Cmul(u), UU[0].Cmul(v))
			<-Z_req
			Z_in <- W.Get()
			Z.Copy(W.Plus(Mul(UU[1], VV[1])))
		}
	}(U, V, Z)
	return Z
}

// Diff erentiate returns the derivative of U.
func (U PS) Diff() PS {
	Z := NewPS()
	go func(U, Z PS) {
		Z_req, Z_dat := Z.Into()
		<-Z_req
		u := U.Get()
		if u.End() == 0 {
			done := false
			for i := 1; !done; i++ {
				u = U.Get()
				if u.End() != 0 {
					done = true
				} else {
					Z_dat <- u.Mul(big.NewRatI(i), u)
					<-Z_req
				}
			}
		}
		Z_dat <- big.Finis
	}(U, Z)
	return Z
}

// Integrate, with const of integration
func (U PS) Integ(c *big.Rat) PS {
	Z := NewPS()
	go func(c *big.Rat, U, Z PS) {
		Z_req, Z_in := Z.Into()
		Z.Put(c)
		done := false
		for i := 1; !done; i++ {
			<-Z_req
			u := U.Get()
			if u.End() == 0 {
				Z_in <- u.Mul(big.NewRat1byI(i), u)
			} else {
				Z_in <- big.Finis
				done = true
			}
		}
	}(c, U, Z)
	return Z
}

// Recip rocal of a power series
//	let U = `u + x*UU`
//	let Z = `z + x*ZZ`
//	`(u+x*UU)*(z+x*ZZ) = 1`
//	`z = 1/u`
//	`u*ZZ + z*UU + x*UU*ZZ = 0`
//
//	ZZ = `-UU * (z + x*ZZ) /u`
func (U PS) Recip() PS {
	Z := NewPS()
	go func(U, Z PS) {
		Z_req, Z_in := Z.Into()
		ZZ := NewPS2()
		<-Z_req
		z := U.Get()
		z.Inv(z)
		Z_in <- z
		var mz *big.Rat
		mz.Neg(z) // minus z `-z`
		ZZ.Split(Mul(U.Cmul(mz), ZZ[0].Shift(z)))
		Z.Copy(ZZ[1])
	}(U, Z)
	return Z
}

// Exp onential of a power series with constant term 0
// (nonzero constant term would make nonrational coefficients)
// BUG: the constant term is simply ignored
//	Z = exp(U)
//	DZ = Z*DU
//	integrate to get Z
func (U PS) Exp() PS {
	ZZ := NewPS2()
	ZZ.Split(Mul(ZZ[0], U.Diff()).Integ(big.One))
	return ZZ[1]
}

// Subst itute V for x in U, where the leading term of V is zero
//	let U = `u + x*UU`
//	let V = `v + x*VV`
//	then S(U,V) = `u + VV*S(V,UU)`
// BUG: a nonzero constant term is ignored
func (U PS) Subst(V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		Z_req, Z_in := Z.Into()
		VV := V.Split()
		<-Z_req
		u := U.Get()
		Z_in <- u
		if u.End() == 0 {
			if VV[0].Get().End() != 0 {
				Z.Put(big.Finis)
			} else {
				Z.Copy(Mul(VV[0], U.Subst(VV[1])))
			}
		}
	}(U, V, Z)
	return Z
}

// MonSubst Monomial Substition: U(c x^n)
// Each Ui is multiplied by `c^i` and followed by n-1 zeros
func (U PS) MonSubst(c0 *big.Rat, n int) PS {
	Z := NewPS()
	go func(U, Z PS, c0 *big.Rat, n int) {
		Z_req, Z_in := Z.Into()
		c := big.One
		var uc *big.Rat // `u * c`
		for {
			<-Z_req
			u := U.Get()
			Z_in <- uc.Mul(u, c)
			c.Mul(c, c0)
			if u.End() != 0 {
				Z_in <- big.Finis
				break
			}
			for i := 1; i < n; i++ {
				<-Z_req
				Z_in <- big.Zero
			}
		}
	}(U, Z, c0, n)
	return Z
}

// ===========================================================================
