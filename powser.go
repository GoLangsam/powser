// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

import (
	"github.com/GoLangsam/powser/big"
	"github.com/GoLangsam/powser/dch"
	"github.com/GoLangsam/powser/dchpair"
)

type PS = *dch.Dch // power series

type PS2 = dchpair.DchPair // pair of power series

var Ones PS
var Twos PS

func NewPS() PS {
	return dch.New()
}

func NewPS2() PS2 {
	return dchpair.DchPair{dch.New(), dch.New()}
}

// GetValS returns a slice with each first value received from the given power series
func GetValS(U ...PS) []*big.Rat {
	pair := dchpair.DchPair{U[0], U[1]}.Get()
	return []*big.Rat{pair[0], pair[1]}
}

// Conventions
// Upper-case for power series.
// Lower-case for rationals.
// Input variables: U,V,...
// Output variables: ...,Y,Z

// Evaln evaluates PS at x=c to n terms in floating point
func Evaln(c *big.Rat, U PS, n int) {
	xn := float64(1)
	x := float64(c.Num()) / float64(c.Den())
	val := float64(0)
	for i := 0; i < n; i++ {
		u := U.Get()
		if u.End() != 0 {
			break
		}
		val = val + x*float64(u.Num())/float64(u.Den())
		xn = xn * x
	}
	print(val, "\n")
}

// Printn prints n terms of a power series
func Printn(U PS, n int) {
	done := false
	for ; !done && n > 0; n-- {
		u := U.Get()
		if u.End() != 0 {
			done = true
		} else {
			u.Pr()
		}
	}
	print(("\n"))
}

// Print one billion terms
func Print(U PS) {
	Printn(U, 1000000000)
}

// eval n terms of power series U at x=c
func eval(c *big.Rat, U PS, n int) *big.Rat {
	if n == 0 {
		return big.Zero
	}
	y := U.Get()
	if y.End() != 0 {
		return big.Zero
	}
	return big.Add(y, big.Mul(c, eval(c, U, n-1)))
}

// Power-series constructors return channels on which power
// series flow.  They start an encapsulated generator that
// puts the terms of the series on the channel.

// Split returns a pair of power series identical to a given power series
func Split(U PS) PS2 {
	UU := NewPS2()
	go UU.Split(U)
	return UU
}

// Arithmetic on power series: each spawns a goroutine

// Add two power series
func Add(U, V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		Z_req, Z_val := Z.Into()
		uv := make([]*big.Rat, 2)
		for {
			<-Z_req
			uv = GetValS(U, V)
			switch uv[0].End() + 2*uv[1].End() {
			case 0:
				Z_val <- big.Add(uv[0], uv[1])
			case 1:
				Z_val <- uv[1]
				Z.Copy(V)
			case 2:
				Z_val <- uv[0]
				Z.Copy(U)
			case 3:
				Z_val <- big.Finis
			}
		}
	}(U, V, Z)
	return Z
}

// Cmul multiplies a power series by a constant
func Cmul(c *big.Rat, U PS) PS {
	Z := NewPS()
	go func(c *big.Rat, U, Z PS) {
		Z_req, Z_val := Z.Into()
		done := false
		for !done {
			<-Z_req
			u := U.Get()
			if u.End() != 0 {
				done = true
			} else {
				Z_val <- big.Mul(c, u)
			}
		}
		Z_val <- big.Finis
	}(c, U, Z)
	return Z
}

// Sub subtracts `V` from `U`
// and returns `U + (-1)*V`
func Sub(U, V PS) PS {
	return Add(U, Cmul(big.MinusOne, V))
}

// Monmul multiplies a power series by the monomial "x^n"
// and returns `x^n * U`.
func Monmul(U PS, n int) PS {
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
func Xmul(U PS) PS {
	return Monmul(U, 1)
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
func Shift(c *big.Rat, U PS) PS {
	Z := NewPS()
	go func(c *big.Rat, U, Z PS) {
		Z.Put(c)
		Z.Copy(U)
	}(c, U, Z)
	return Z
}

// simple pole at 1: 1/(1-x) = 1 1 1 1 1 ...

// Poly converts coefficients, constant term first
// to a (finite) power series
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
		Z_req, Z_val := Z.Into()
		<-Z_req
		uv := GetValS(U, V)
		if uv[0].End() != 0 || uv[1].End() != 0 {
			Z_val <- big.Finis
		} else {
			Z_val <- big.Mul(uv[0], uv[1])
			UU := Split(U)
			VV := Split(V)
			W := Add(Cmul(uv[0], VV[0]), Cmul(uv[1], UU[0]))
			<-Z_req
			Z_val <- W.Get()
			Z.Copy(Add(W, Mul(UU[1], VV[1])))
		}
	}(U, V, Z)
	return Z
}

// Diff erentiate
func Diff(U PS) PS {
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
					Z_dat <- big.Mul(big.NewRatI(i), u)
					<-Z_req
				}
			}
		}
		Z_dat <- big.Finis
	}(U, Z)
	return Z
}

// Integrate, with const of integration
func Integ(c *big.Rat, U PS) PS {
	Z := NewPS()
	go func(c *big.Rat, U, Z PS) {
		Z_req, Z_val := Z.Into()
		Z.Put(c)
		done := false
		for i := 1; !done; i++ {
			<-Z_req
			u := U.Get()
			if u.End() == 0 {
				Z_val <- big.Mul(big.NewRat1byI(i), u)
			} else {
				Z_val <- big.Finis
				done = true
			}
		}
	}(c, U, Z)
	return Z
}

// Binom ial theorem
// and returns `(1+x)^c`
func Binom(c *big.Rat) PS {
	Z := NewPS()
	go func(c *big.Rat, Z PS) {
		n := 1
		t := big.One
		for c.Num() != 0 {
			Z.Put(t)
			t = big.Mul(big.Mul(t, c), big.NewRat1byI(n))
			c = big.Sub(c, big.One)
			n++
		}
		Z.Put(big.Finis)
	}(c, Z)
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
func Recip(U PS) PS {
	Z := NewPS()
	go func(U, Z PS) {
		req, dat := Z.Into()
		ZZ := NewPS2()
		<-req
		z := big.Inv(U.Get())
		dat <- z
		ZZ.Split(Mul(Cmul(big.Neg(z), U), Shift(z, ZZ[0])))
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
func Exp(U PS) PS {
	ZZ := NewPS2()
	ZZ.Split(Integ(big.One, Mul(ZZ[0], Diff(U))))
	return ZZ[1]
}

// Subst itute V for x in U, where the leading term of V is zero
//	let U = `u + x*UU`
//	let V = `v + x*VV`
//	then S(U,V) = `u + VV*S(V,UU)`
// BUG: a nonzero constant term is ignored
func Subst(U, V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		req, dat := Z.Into()
		VV := Split(V)
		<-req
		u := U.Get()
		dat <- u
		if u.End() == 0 {
			if VV[0].Get().End() != 0 {
				Z.Put(big.Finis)
			} else {
				Z.Copy(Mul(VV[0], Subst(U, VV[1])))
			}
		}
	}(U, V, Z)
	return Z
}

// MonSubst Monomial Substition: U(c x^n)
// Each Ui is multiplied by `c^i` and followed by n-1 zeros
func MonSubst(U PS, c0 *big.Rat, n int) PS {
	Z := NewPS()
	go func(U, Z PS, c0 *big.Rat, n int) {
		req, dat := Z.Into()
		c := big.One
		for {
			<-req
			u := U.Get()
			dat <- big.Mul(u, c)
			c = big.Mul(c, c0)
			if u.End() != 0 {
				dat <- big.Finis
				break
			}
			for i := 1; i < n; i++ {
				<-req
				dat <- big.Zero
			}
		}
	}(U, Z, c0, n)
	return Z
}
