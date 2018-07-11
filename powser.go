// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

import (
	"github.com/GoLangsam/powser/rat"
)

type PS *dch    // power series
type PS2 *[2]PS // pair of power series

var Ones PS
var Twos PS

func mkPS() *dch {
	return mkdch()
}

func mkPS2() *dch2 {
	return mkdch2()
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
		u := get(U)
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
		u := get(U)
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
	y := get(U)
	if y.End() != 0 {
		return rat.Zero
	}
	return rat.Add(y, rat.Mul(c, eval(c, U, n-1)))
}

// Power-series constructors return channels on which power
// series flow.  They start an encapsulated generator that
// puts the terms of the series on the channel.

// Make a pair of power series identical to a given power series

func Split(U PS) *dch2 {
	UU := mkdch2()
	go split(U, UU)
	return UU
}

// Add two power series
func Add(U, V PS) PS {
	Z := mkPS()
	go func(U, V, Z PS) {
		var uv [2]*rat.Rat
		for {
			<-Z.req
			uv = get2([2]*dch{U, V})
			switch uv[0].End() + 2*uv[1].End() {
			case 0:
				Z.dat <- rat.Add(uv[0], uv[1])
			case 1:
				Z.dat <- uv[1]
				copy(V, Z)
			case 2:
				Z.dat <- uv[0]
				copy(U, Z)
			case 3:
				Z.dat <- rat.Finis
			}
		}
	}(U, V, Z)
	return Z
}

// Multiply a power series by a constant
func Cmul(c *rat.Rat, U PS) PS {
	Z := mkPS()
	go func(c *rat.Rat, U, Z PS) {
		done := false
		for !done {
			<-Z.req
			u := get(U)
			if u.End() != 0 {
				done = true
			} else {
				Z.dat <- rat.Mul(c, u)
			}
		}
		Z.dat <- rat.Finis
	}(c, U, Z)
	return Z
}

// Subtract

func Sub(U, V PS) PS {
	return Add(U, Cmul(rat.MinusOne, V))
}

// Multiply a power series by the monomial x^n

func Monmul(U PS, n int) PS {
	Z := mkPS()
	go func(n int, U PS, Z PS) {
		for ; n > 0; n-- {
			put(rat.Zero, Z)
		}
		copy(U, Z)
	}(n, U, Z)
	return Z
}

// Multiply by x

func Xmul(U PS) PS {
	return Monmul(U, 1)
}

func Rep(c *rat.Rat) PS {
	Z := mkPS()
	go repeat(c, Z)
	return Z
}

// Monomial c*x^n

func Mon(c *rat.Rat, n int) PS {
	Z := mkPS()
	go func(c *rat.Rat, n int, Z PS) {
		if c.Num() != 0 {
			for ; n > 0; n = n - 1 {
				put(rat.Zero, Z)
			}
			put(c, Z)
		}
		put(rat.Finis, Z)
	}(c, n, Z)
	return Z
}

func Shift(c *rat.Rat, U PS) PS {
	Z := mkPS()
	go func(c *rat.Rat, U, Z PS) {
		put(c, Z)
		copy(U, Z)
	}(c, U, Z)
	return Z
}

// simple pole at 1: 1/(1-x) = 1 1 1 1 1 ...

// Convert array of coefficients, constant term first
// to a (finite) power series

/*
func Poly(a [] *rat.Rat) PS{
	Z:=mkPS()
	begin func(a [] *rat.Rat, Z PS){
		j:=0
		done:=0
		for j=len(a); !done&&j>0; j=j-1)
			if(a[j-1].num!=0) done=1
		i:=0
		for(; i<j; i=i+1) put(a[i],Z)
		put(rat.Finis,Z)
	}()
	return Z
}
*/

// Multiply. The algorithm is
//	let U = u + x*UU
//	let V = v + x*VV
//	then UV = u*v + x*(u*VV+v*UU) + x*x*UU*VV

func Mul(U, V PS) PS {
	Z := mkPS()
	go func(U, V, Z PS) {
		<-Z.req
		uv := get2([2]*dch{U, V})
		if uv[0].End() != 0 || uv[1].End() != 0 {
			Z.dat <- rat.Finis
		} else {
			Z.dat <- rat.Mul(uv[0], uv[1])
			UU := Split(U)
			VV := Split(V)
			W := Add(Cmul(uv[0], VV[0]), Cmul(uv[1], UU[0]))
			<-Z.req
			Z.dat <- get(W)
			copy(Add(W, Mul(UU[1], VV[1])), Z)
		}
	}(U, V, Z)
	return Z
}

// Differentiate

func Diff(U PS) PS {
	Z := mkPS()
	go func(U, Z PS) {
		<-Z.req
		u := get(U)
		if u.End() == 0 {
			done := false
			for i := 1; !done; i++ {
				u = get(U)
				if u.End() != 0 {
					done = true
				} else {
					Z.dat <- rat.Mul(rat.ItoR(int64(i)), u)
					<-Z.req
				}
			}
		}
		Z.dat <- rat.Finis
	}(U, Z)
	return Z
}

// Integrate, with const of integration
func Integ(c *rat.Rat, U PS) PS {
	Z := mkPS()
	go func(c *rat.Rat, U, Z PS) {
		put(c, Z)
		done := false
		for i := 1; !done; i++ {
			<-Z.req
			u := get(U)
			if u.End() != 0 {
				done = true
			}
			Z.dat <- rat.Mul(rat.I2toR(1, int64(i)), u)
		}
		Z.dat <- rat.Finis
	}(c, U, Z)
	return Z
}

// Binomial theorem (1+x)^c

func Binom(c *rat.Rat) PS {
	Z := mkPS()
	go func(c *rat.Rat, Z PS) {
		n := 1
		t := rat.ItoR(1)
		for c.Num() != 0 {
			put(t, Z)
			t = rat.Mul(rat.Mul(t, c), rat.I2toR(1, int64(n)))
			c = rat.Sub(c, rat.One)
			n++
		}
		put(rat.Finis, Z)
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
	Z := mkPS()
	go func(U, Z PS) {
		ZZ := mkPS2()
		<-Z.req
		z := rat.Inv(get(U))
		Z.dat <- z
		split(Mul(Cmul(rat.Neg(z), U), Shift(z, ZZ[0])), ZZ)
		copy(ZZ[1], Z)
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
	ZZ := mkPS2()
	split(Integ(rat.One, Mul(ZZ[0], Diff(U))), ZZ)
	return ZZ[1]
}

// Substitute V for x in U, where the leading term of V is zero
//	let U = u + x*UU
//	let V = v + x*VV
//	then S(U,V) = u + VV*S(V,UU)
// bug: a nonzero constant term is ignored

func Subst(U, V PS) PS {
	Z := mkPS()
	go func(U, V, Z PS) {
		VV := Split(V)
		<-Z.req
		u := get(U)
		Z.dat <- u
		if u.End() == 0 {
			if get(VV[0]).End() != 0 {
				put(rat.Finis, Z)
			} else {
				copy(Mul(VV[0], Subst(U, VV[1])), Z)
			}
		}
	}(U, V, Z)
	return Z
}

// Monomial Substition: U(c x^n)
// Each Ui is multiplied by c^i and followed by n-1 zeros

func MonSubst(U PS, c0 *rat.Rat, n int) PS {
	Z := mkPS()
	go func(U, Z PS, c0 *rat.Rat, n int) {
		c := rat.One
		for {
			<-Z.req
			u := get(U)
			Z.dat <- rat.Mul(u, c)
			c = rat.Mul(c, c0)
			if u.End() != 0 {
				Z.dat <- rat.Finis
				break
			}
			for i := 1; i < n; i++ {
				<-Z.req
				Z.dat <- rat.Zero
			}
		}
	}(U, Z, c0, n)
	return Z
}
