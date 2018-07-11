// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================
// Arithmetic on power series: each spawns a goroutine

// Conventions
// Upper-case for power series.
// Lower-case for coefficients.
// Input variables: U,V,...
// Output variables: ...,Y,Z

// Add two power series.
func Add(U, V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		Z_req, Z_in := Z.Into()
		defer Z.Close()

		var u, v Coefficient
		var uok, vok bool
		for {
			<-Z_req
			u, v, uok, vok = GetVal2(U, V)
			switch { // fini(u) + 2*fini(v) {
			case uok && vok:
				Z_in <- u.Add(u, v)
			case uok:
				Z_in <- u
				Z.Append(U)
			case vok:
				Z_in <- v
				Z.Append(V)
			default:
				return
			}
		}
	}(U, V, Z)
	return Z
}

// Plus adds powerseries to `U` and returns the sum.
// Tail-recursion is used to achieve this.
func (U PS) Plus(V ...PS) PS {
	switch len(V) {
	case 0:
		return U
	case 1:
		return Add(U, V[0])
	default:
		return Add(U, V[0]).Plus(V[1:]...)
	}
}

// Cmul multiplies a power series by a constant
func (U PS) Cmul(c Coefficient) PS {
	Z := NewPS()
	go func(c Coefficient, U, Z PS) {
		Z_req, Z_in := Z.Into()
		defer Z.Close()

		var u Coefficient
		var ok bool
		for {
			<-Z_req
			if u, ok = U.Get(); !ok {
				return
			}
			Z_in <- u.Mul(c, u)
		}
	}(c, U, Z)
	return Z
}

// Minus subtracts `V` from `U`
// and returns `U + (-1)*V`
func (U PS) Minus(V PS) PS {
	return U.Plus(V.Cmul(aMinusOne))
}

// Less subtracts powerseries from `U` and returns the difference.
// Tail-recursion is used to achieve this.
func (U PS) Less(V ...PS) PS {
	switch len(V) {
	case 0:
		return U
	case 1:
		return U.Minus(V[0])
	default:
		return U.Minus(V[0]).Less(V[1:]...)
	}
}

// Monmul multiplies `U` by the monomial "x^n"
// and returns `x^n * U`.
func (U PS) Monmul(n int) PS {
	Z := NewPS()
	go func(n int, U PS, Z PS) {
		for ; n > 0; n-- {
			Z.Put(aZero)
		}
		Z.Append(U)
	}(n, U, Z)
	return Z
}

// Xmul multiplies a power series by x, (by the monomial "x^1")
// and returns `x * U`.
func (U PS) Xmul() PS {
	return U.Monmul(1)
}

// Shift
func (U PS) Shift(c Coefficient) PS {
	Z := NewPS()
	go func(c Coefficient, U, Z PS) {
		Z.Put(c)
		Z.Append(U)
	}(c, U, Z)
	return Z
}

// simple pole at 1: 1/(1-x) = 1 1 1 1 1 ...

// Mul multiplies. The algorithm is
//	let U = `u + x*UU`
//	let V = `v + x*VV`
//	then UV = `u*v + x*(u*VV+v*UU) + x*x*UU*VV`
func Mul(U, V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		Z_req, Z_in := Z.Into()
		defer Z.Close()

		<-Z_req
		u, v, uok, vok := GetVal2(U, V)
		if !uok || !vok {
			return
		}
		var c Coefficient
		c.Mul(u, v)
		Z_in <- c
		UU := U.split()
		VV := V.split()
		W := Add(VV[0].Cmul(u), UU[0].Cmul(v))
		<-Z_req
		c, _ = W.Get()
		Z_in <- c
		Z.Append(W.Plus(Mul(UU[1], VV[1])))

	}(U, V, Z)
	return Z
}

// Times multiplies powerseries to `U` and returns the total product.
// Tail-recursion is used to achieve this.
func (U PS) Times(V ...PS) PS {
	switch len(V) {
	case 0:
		return U
	case 1:
		return Mul(U, V[0])
	default:
		return Mul(U, V[0]).Times(V[1:]...)
	}
}

// Diff erentiate returns the derivative of U.
func (U PS) Diff() PS {
	Z := NewPS()
	go func(U, Z PS) {
		Z_req, Z_in := Z.Into()
		defer Z.Close()

		<-Z_req

		u, ok := U.Get()
		if !ok {
			return
		}

		for i := 1; ; i++ {
			if u, ok = U.Get(); !ok {
				return
			}
			Z_in <- u.Mul(ratI(i), u)
			<-Z_req
		}

	}(U, Z)
	return Z
}

// Integrate, with const of integration
func (U PS) Integ(c Coefficient) PS {
	Z := NewPS()
	go func(c Coefficient, U, Z PS) {
		Z_req, Z_in := Z.Into()
		defer Z.Close()

		Z.Put(c)

		var u Coefficient
		var ok bool
		for i := 1; ; i++ {
			<-Z_req
			if u, ok = U.Get(); !ok {
				return
			}
			Z_in <- u.Mul(rat1byI(i), u)
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
//	ZZ = `1/u * -UU * (z + x*ZZ)`
func (U PS) Recip() PS {
	Z := NewPS()
	go func(U, Z PS) {
		Z_req, Z_in := Z.Into()
		defer Z.Close()

		<-Z_req
		z, ok := U.Get()
		if !ok {
			return
		}

		z.Inv(z)
		Z_in <- z
		var mz Coefficient
		mz.Neg(z) // minus z `-z`
		ZZ := NewPS2()
		ZZ.split(Mul(U.Cmul(mz), ZZ[0].Shift(z)))
		Z.Append(ZZ[1])

	}(U, Z)
	return Z
}

// Exp onential of a power series with constant term 0
// (nonzero constant term would make nonrational coefficients)
//	Z = exp(U)
//	DZ = Z*DU
//	integrate to get Z
// Note: The constant term is simply ignored.
func (U PS) Exp() PS {
	ZZ := NewPS2()
	ZZ.split(Mul(ZZ[0], U.Diff()).Integ(aOne))
	return ZZ[1]
}

// Subst itute V for x in U, where the leading term of V is zero
//	let U = `u + x*UU`
//	let V = `v + x*VV`
//	then Subst(U,V) = `u + VV * Subst(V,UU)`
// Note: Any nonzero constant term of `V` is ignored.
func (U PS) Subst(V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		Z_req, Z_in := Z.Into()
		defer Z.Close()

		var u Coefficient
		var ok bool

		VV := V.split()
		<-Z_req

		if u, ok = U.Get(); !ok {
			return
		}

		Z_in <- u

		if _, ok = VV[0].Get(); !ok {
			return // Note: Any nonzero constant term of `V` is ignored.
		}

		Z.Append(Mul(VV[0], U.Subst(VV[1])))

	}(U, V, Z)
	return Z
}

// MonSubst Monomial Substition: U(c x^n)
// Each Ui is multiplied by `c^i` and followed by n-1 zeros.
func (U PS) MonSubst(c0 Coefficient, n int) PS {
	Z := NewPS()
	go func(U, Z PS, c0 Coefficient, n int) {
		Z_req, Z_in := Z.Into()
		defer Z.Close()

		var u Coefficient
		var ok bool
		c := aOne
		var uc Coefficient // `u * c`
		for {
			<-Z_req
			if u, ok = U.Get(); !ok {
				return
			}

			Z_in <- uc.Mul(u, c)
			c.Mul(c, c0)

			for i := 1; i < n; i++ {
				<-Z_req
				Z_in <- aZero
			}
		}
	}(U, Z, c0, n)
	return Z
}

// ===========================================================================
