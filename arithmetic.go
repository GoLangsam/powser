// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================
// Arithmetic on power series: each spawns a goroutine

// Conventions
// Upper-case for power series.
// Lower-case for coefficients.
// Input variables:  From U,V,...
// Output variables: Into ...,Y,Z

// Add two power series.
func Add(U, V PS) PS {
	Z := U.New()
	go func(U, V, Z PS) {
		defer func() {
			U.Drop()
			V.Drop()
			Z.Close()
		}()

		for Z.Next() {
			u, okU, v, okV := U.GetWith(V)
			if !okU && !okV {
				return
			}
			Z.Send(aC().Add(u, v)) // `u + v`
		}
	}(U, V, Z)
	return Z
}

// Plus adds powerseries to `U`
// and returns the sum.
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

// Minus subtracts `V` from `U`
// and returns `U + (-1)*V`
func (U PS) Minus(V PS) PS {
	return U.Plus(V.CMul(aMinusOne()))
}

// Less subtracts powerseries from `U`
// and returns the difference.
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

// CMul multiplies `U` by a constant `c`
// and returns `c*U`.
func (U PS) CMul(c Coefficient) PS {
	Z := U.New()
	go func(c Coefficient, U, Z PS) {
		for Z.SendCfnFrom(U, cMul(c)) { // `c * u`
		}
	}(c, U, Z)
	return Z
}

// MonMul multiplies `U` by the monomial "x^n"
// and returns `x^n * U`.
func (U PS) MonMul(n int) PS {
	Z := U.New()
	go func(n int, U PS, Z PS) {
		for ; n > 0; n-- {
			Z.Put(aZero())
		}
		Z.Append(U)
	}(n, U, Z)
	return Z
}

// XMul multiplies `U` by `x`
// (by the monomial "x^1")
// and returns `x * U`.
func (U PS) XMul() PS {
	return U.MonMul(1)
}

// Shift returns `c + x*U`
func (U PS) Shift(c Coefficient) PS {
	Z := U.New()
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
	Z := U.New()
	go func(U, V, Z PS) {
		var u, v Coefficient
		var next, okU, okV bool

		if next = Z.Next(); next {
			u, okU, v, okV = U.GetWith(V)
		}

		if !next || !okU || !okV { // Z.Dropped or U or V.Closed
			U.Drop()
			V.Drop()
			Z.Close()
			return
		}

		Z.Send(aC().Mul(u, v))                 // `u*v`
		UU, VV := U.Split(), V.Split()         // `UU`, `VV`
		W := Add(VV[0].CMul(u), UU[0].CMul(v)) // `u*VV + v*UU`
		if Z.SendCfnFrom(W, cSame()) {         // ` + x*(u*VV+v*UU)`
			Z.Append(W.Plus(Mul(UU[1], VV[1]))) // `+ x*x*UU*VV`
		}
	}(U, V, Z)
	return Z
}

// Times multiplies powerseries to `U`
// and returns the total product.
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

// Deriv differentiates `U`
// and returns the derivative.
func (U PS) Deriv() PS {
	Z := U.New()
	go func(U, Z PS) {
		if _, ok := Z.GetNextFrom(U); !ok {
			return
		}
		// constant term: drop
		// Thus: we must Z.Send() before another Z.Next()
		// and may not use an Obtain-loop and have to cleanup ourselfs

		for i := 1; ; i++ {
			if u, ok := U.Get(); ok {
				Z.Send(cRatIby1(i)(u)) // `u * i`
				if !Z.Next() {
					break
				}
			} else {
				break
			}
		}
		Z.Close()
		U.Drop()
	}(U, Z)
	return Z
}

// Integrate, with const of integration.
func (U PS) Integ(c Coefficient) PS {
	Z := U.New()
	go func(c Coefficient, U, Z PS) {
		Z.Put(c)

		i := 1
		for Z.SendCfnFrom(U, cRat1byI(i)) { // `u * 1/i`
			i++
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
//	ZZ = `1/u * (-z*UU + x*UU*ZZ)`
func (U PS) Recip() PS {
	Z := U.New()
	go func(U, Z PS) {
		u, ok := Z.GetNextFrom(U)
		if !ok {
			return
		}
		Z.Send(aC().Inv(u)) // `1/u`
		mu := aC().Neg(u)   // `-z` minus z
		ZZ := U.newPair()
		ZZ.Split(Mul(U.CMul(mu), ZZ[0].Shift(u)))
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
	ZZ := U.newPair()
	ZZ.Split(Mul(ZZ[0], U.Deriv()).Integ(aOne()))
	return ZZ[1]
}

// Subst itute V for x in U, where the leading term of V is zero
//	let U = `u + x*UU`
//	let V = `v + x*VV`
//	then Subst(U,V) = `u + VV * Subst(V,UU)`
// Note: Any nonzero constant term of `V` is ignored.
func (U PS) Subst(V PS) PS {
	Z := U.New()
	go func(U, V, Z PS) {
		if !Z.SendCfnFrom(U, cSame()) {
			V.Drop()
			return
		}

		VV := V.Split()
		VV[0].Get() // Note: Any nonzero constant term of `V` is ignored.
		Z.Append(Mul(VV[0], U.Subst(VV[1])))

	}(U, V, Z)
	return Z
}

// MonSubst Monomial Substition: `U(c x^n)`
//
// Each Ui is multiplied by `c^i` and followed by n-1 zeros.
func (U PS) MonSubst(c0 Coefficient, n int) PS {
	Z := U.New()
	go func(U, Z PS, c0 Coefficient, n int) {
		c := aOne()
		for Z.SendCfnFrom(U, cMul(c)) { // `c * u`
			for i := 1; i < n; i++ {
				if !Z.Put(aZero()) { // n-1 zeros
					Z.Close()
					U.Drop()
					return
				}
			}
			c.Mul(c, c0) // `c = c * c0 = c^i`
		}
	}(U, Z, c0, n)
	return Z
}

// ===========================================================================
