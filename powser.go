// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

import (
	"github.com/GoLangsam/powser/dch"
	"github.com/GoLangsam/powser/rat"
)

type PS = *dch.Dch // power series

type PS2 = dch.DchPair // pair of power series

var Ones PS
var Twos PS

func NewPS() PS {
	return dch.New()
}

func NewPS2() PS2 {
	return dch.NewPair()
}

// Conventions
// Upper-case for power series.
// Lower-case for rationals.
// Input variables: U,V,...
// Output variables: ...,Y,Z

// print eval in floating point of PS at x=c to n terms
func Evaln(c *rat.Rat, U PS, n int) {
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

// Print n terms of a power series
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

func Print(U PS) {
	Printn(U, 1000000000)
}

// Evaluate n terms of power series U at x=c
func eval(c *rat.Rat, U PS, n int) *rat.Rat {
	if n == 0 {
		return rat.Zero
	}
	y := U.Get()
	if y.End() != 0 {
		return rat.Zero
	}
	return rat.Add(y, rat.Mul(c, eval(c, U, n-1)))
}

// Power-series constructors return channels on which power
// series flow.  They start an encapsulated generator that
// puts the terms of the series on the channel.

// Make a pair of power series identical to a given power series

func Split(U PS) [2]*dch.Dch {
	UU := NewPS2()
	go UU.Split(U)
	return UU
}

// Add two power series
func Add(U, V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		var uv [2]*rat.Rat
		for {
			<-Z.Req()
			uv = dch.Get2(dch.DchPair{U, V})
			switch uv[0].End() + 2*uv[1].End() {
			case 0:
				Z.Dat() <- rat.Add(uv[0], uv[1])
			case 1:
				Z.Dat() <- uv[1]
				Z.Copy(V)
			case 2:
				Z.Dat() <- uv[0]
				Z.Copy(U)
			case 3:
				Z.Dat() <- rat.Finis
			}
		}
	}(U, V, Z)
	return Z
}

// Multiply a power series by a constant
func Cmul(c *rat.Rat, U PS) PS {
	Z := NewPS()
	go func(c *rat.Rat, U, Z PS) {
		done := false
		for !done {
			<-Z.Req()
			u := U.Get()
			if u.End() != 0 {
				done = true
			} else {
				Z.Dat() <- rat.Mul(c, u)
			}
		}
		Z.Dat() <- rat.Finis
	}(c, U, Z)
	return Z
}

// Subtract

func Sub(U, V PS) PS {
	return Add(U, Cmul(rat.MinusOne, V))
}

// Multiply a power series by the monomial x^n

func Monmul(U PS, n int) PS {
	Z := NewPS()
	go func(n int, U PS, Z PS) {
		for ; n > 0; n-- {
			Z.Put(rat.Zero)
		}
		Z.Copy(U)
	}(n, U, Z)
	return Z
}

// Multiply by x

func Xmul(U PS) PS {
	return Monmul(U, 1)
}

func Rep(c *rat.Rat) PS {
	Z := NewPS()
	go Z.Repeat(c)
	return Z
}

// Monomial c*x^n

func Mon(c *rat.Rat, n int) PS {
	Z := NewPS()
	go func(c *rat.Rat, n int, Z PS) {
		if c.Num() != 0 {
			for ; n > 0; n = n - 1 {
				Z.Put(rat.Zero)
			}
			Z.Put(c)
		}
		Z.Put(rat.Finis)
	}(c, n, Z)
	return Z
}

func Shift(c *rat.Rat, U PS) PS {
	Z := NewPS()
	go func(c *rat.Rat, U, Z PS) {
		Z.Put(c)
		Z.Copy(U)
	}(c, U, Z)
	return Z
}

// simple pole at 1: 1/(1-x) = 1 1 1 1 1 ...

// Convert array of coefficients, constant term first
// to a (finite) power series

/*
func Poly(a [] *rat.Rat) PS{
	Z:=NewPS()
	begin func(a [] *rat.Rat, Z PS){
		j:=0
		done:=0
		for j=len(a); !done&&j>0; j=j-1)
			if(a[j-1].num!=0) done=1
		i:=0
		for(; i<j; i=i+1) Z.Put(a[i])
		Z.Put(rat.Finis)
	}()
	return Z
}
*/

// Multiply. The algorithm is
//	let U = u + x*UU
//	let V = v + x*VV
//	then UV = u*v + x*(u*VV+v*UU) + x*x*UU*VV

func Mul(U, V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		<-Z.Req()
		uv := dch.Get2([2]*dch.Dch{U, V})
		if uv[0].End() != 0 || uv[1].End() != 0 {
			Z.Dat() <- rat.Finis
		} else {
			Z.Dat() <- rat.Mul(uv[0], uv[1])
			UU := Split(U)
			VV := Split(V)
			W := Add(Cmul(uv[0], VV[0]), Cmul(uv[1], UU[0]))
			<-Z.Req()
			Z.Dat() <- W.Get()
			Z.Copy(Add(W, Mul(UU[1], VV[1])))
		}
	}(U, V, Z)
	return Z
}

// Differentiate

func Diff(U PS) PS {
	Z := NewPS()
	go func(U, Z PS) {
		<-Z.Req()
		u := U.Get()
		if u.End() == 0 {
			done := false
			for i := 1; !done; i++ {
				u = U.Get()
				if u.End() != 0 {
					done = true
				} else {
					Z.Dat() <- rat.Mul(rat.NewRat(int64(i), 1), u)
					<-Z.Req()
				}
			}
		}
		Z.Dat() <- rat.Finis
	}(U, Z)
	return Z
}

// Integrate, with const of integration
func Integ(c *rat.Rat, U PS) PS {
	Z := NewPS()
	go func(c *rat.Rat, U, Z PS) {
		Z.Put(c)
		done := false
		for i := 1; !done; i++ {
			<-Z.Req()
			u := U.Get()
			if u.End() != 0 {
				done = true
			}
			Z.Dat() <- rat.Mul(rat.NewRat(1, int64(i)), u)
		}
		Z.Dat() <- rat.Finis
	}(c, U, Z)
	return Z
}

// Binomial theorem (1+x)^c

func Binom(c *rat.Rat) PS {
	Z := NewPS()
	go func(c *rat.Rat, Z PS) {
		n := 1
		t := rat.One
		for c.Num() != 0 {
			Z.Put(t)
			t = rat.Mul(rat.Mul(t, c), rat.NewRat(1, int64(n)))
			c = rat.Sub(c, rat.One)
			n++
		}
		Z.Put(rat.Finis)
	}(c, Z)
	return Z
}

// Reciprocal of a power series
//	let U = u + x*UU
//	let Z = z + x*ZZ
//	(u+x*UU)*(z+x*ZZ) = 1
//	z = 1/u
//	u*ZZ + z*UU +x*UU*ZZ = 0
//	ZZ = -UU*(z+x*ZZ)/u

func Recip(U PS) PS {
	Z := NewPS()
	go func(U, Z PS) {
		ZZ := NewPS2()
		<-Z.Req()
		z := rat.Inv(U.Get())
		Z.Dat() <- z
		ZZ.Split(Mul(Cmul(rat.Neg(z), U), Shift(z, ZZ[0])))
		Z.Copy(ZZ[1])
	}(U, Z)
	return Z
}

// Exponential of a power series with constant term 0
// (nonzero constant term would make nonrational coefficients)
// bug: the constant term is simply ignored
//	Z = exp(U)
//	DZ = Z*DU
//	integrate to get Z

func Exp(U PS) PS {
	ZZ := NewPS2()
	ZZ.Split(Integ(rat.One, Mul(ZZ[0], Diff(U))))
	return ZZ[1]
}

// Substitute V for x in U, where the leading term of V is zero
//	let U = u + x*UU
//	let V = v + x*VV
//	then S(U,V) = u + VV*S(V,UU)
// bug: a nonzero constant term is ignored

func Subst(U, V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		VV := Split(V)
		<-Z.Req()
		u := U.Get()
		Z.Dat() <- u
		if u.End() == 0 {
			if VV[0].Get().End() != 0 {
				Z.Put(rat.Finis)
			} else {
				Z.Copy(Mul(VV[0], Subst(U, VV[1])))
			}
		}
	}(U, V, Z)
	return Z
}

// Monomial Substition: U(c x^n)
// Each Ui is multiplied by c^i and followed by n-1 zeros

func MonSubst(U PS, c0 *rat.Rat, n int) PS {
	Z := NewPS()
	go func(U, Z PS, c0 *rat.Rat, n int) {
		c := rat.One
		for {
			<-Z.Req()
			u := U.Get()
			Z.Dat() <- rat.Mul(u, c)
			c = rat.Mul(c, c0)
			if u.End() != 0 {
				Z.Dat() <- rat.Finis
				break
			}
			for i := 1; i < n; i++ {
				<-Z.Req()
				Z.Dat() <- rat.Zero
			}
		}
	}(U, Z, c0, n)
	return Z
}
