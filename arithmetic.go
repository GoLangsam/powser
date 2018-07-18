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

// add two power series.
func add(U, V PS) PS {
	Z := U.new()
	go func(Z PS, U, V PS) {
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
	}(Z, U, V)
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
		return add(U, V[0])
	default:
		return add(U, V[0]).Plus(V[1:]...)
	}
}

// minus subtracts `V` from `U`
// and returns `U + (-1)*V`
func (U PS) minus(V PS) PS {
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
		return U.minus(V[0])
	default:
		return U.minus(V[0]).Less(V[1:]...)
	}
}

// CMul multiplies `U` by a constant `c`
// and returns `c*U`.
func (U PS) CMul(c Coefficient) PS {
	Z := U.new()
	go func(Z PS, U PS, c Coefficient) {
		for Z.SendCfnFrom(U, cMul(c)) { // `c * u`
		}
	}(Z, U, c)
	return Z
}

// MonMul multiplies `U` by the monomial "x^n"
// and returns `x^n * U`.
// If `n` is not positive, zero (an empty closed power series) is returned.
func (U PS) MonMul(n int) PS {
	Z := U.new()

	if !(n > 0) {
		Z.Close()
		U.Drop()
		return Z
	}

	go func(Z PS, U PS, n int) {
		for ; n > 0; n-- {
			if !Z.Put(aZero()) {
				U.Drop()
				return
			}
		}
		Z.Append(U)
	}(Z, U, n)
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
	Z := U.new()
	go func(Z PS, U PS, c Coefficient) {
		if !Z.Put(c) {
			U.Drop()
			return
		}
		Z.Append(U)
	}(Z, U, c)
	return Z
}

// mul multiplies. The algorithm is:
//	let U = `u + x*UU`
//	let V = `v + x*VV`
//	then UV = `u*v + x*(u*VV+v*UU) + x*x*UU*VV`
func mul(U, V PS) PS {
	Z := U.new()
	go func(Z PS, U, V PS) {
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
		W := VV[0].CMul(u).Plus(UU[0].CMul(v)) // `u*VV + v*UU`
		if Z.SendCfnFrom(W, cSame()) {         // ` + x*(u*VV+v*UU)`
			Z.Append(W.Plus(UU[1].Times(VV[1]))) // `+ x*x*UU*VV` - recurse
		}
	}(Z, U, V)
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
		return mul(U, V[0])
	default:
		return mul(U, V[0]).Times(V[1:]...)
	}
}

// Deriv differentiates `U`
// and returns the derivative.
func (U PS) Deriv() PS {
	Z := U.new()
	go func(Z PS, U PS) {
		u, ok := Z.NextGetFrom(U)
		if !ok {
			return
		}
		// constant term: drop
		// Thus: we must Z.Send() before another Z.Next(),
		for i := 1; ok; i++ {
			if u, ok = U.Get(); ok {
				Z.Send(cRatIby1(i)(u)) // `u * i`
				ok = Z.Next()
			}
		}
		Z.Close()
		U.Drop()
	}(Z, U)
	return Z
}

// Integ integrates `U` with `c` as constant of integration.
func (U PS) Integ(c Coefficient) PS {
	Z := U.new()
	go func(Z PS, U PS, c Coefficient) {
		if !Z.Put(c) {
			U.Drop()
			return
		}

		i := 1
		for Z.SendCfnFrom(U, cRat1byI(i)) { // `u * 1/i`
			i++
		}
	}(Z, U, c)
	return Z
}

// Recip rocal of a power series. The algorithm is:
//	let U = `u + x*UU`
//	let Z = `z + x*ZZ`
//	`(u+x*UU)*(z+x*ZZ) = 1`
//	`z = 1/u`
//	`u*ZZ + z*UU + x*UU*ZZ = 0`
//
//	ZZ = `1/u * -UU * (z + x*ZZ)`
//	ZZ = `1/u * (-z*UU + x*UU*ZZ)`
func (U PS) Recip() PS {
	Z := U.new()
	go func(Z PS, U PS) {
		u, ok := Z.NextGetFrom(U)
		if !ok {
			return
		}
		Z.Send(aC().Inv(u)) // `1/u`
		mu := aC().Neg(u)   // `-z` minus z
		ZZ := U.newPair()
		ZZ.Split(U.CMul(mu).Times(ZZ[0].Shift(u)))
		Z.Append(ZZ[1])
	}(Z, U)
	return Z
}

// Exp onentiation of a power series (with constant term equal zero):
//	Z = exp(U)
//	DZ = Z*DU
//	integrate to get Z
// Note: The constant term is simply ignored as
// any nonzero constant term would imply nonrational coefficients.
func (U PS) Exp() PS {
	ZZ := U.newPair()
	ZZ.Split(ZZ[0].Times(U.Deriv()).Integ(aOne()))
	return ZZ[1]
}

// Subst itute V for x in U, where the constant term of V is zero:
//	let U = `u + x*UU`
//	let V = `v + x*VV`
//	then U.Subst(V) = `u + VV * U.Subst(VV)`
// Note: Any nonzero constant term of `V` is simply ignored.
func (U PS) Subst(V PS) PS {
	Z := U.new()
	go func(Z PS, U, V PS) {
		if !Z.SendCfnFrom(U, cSame()) {
			return
		}

		VV := V.Split()
		VV[0].Get() // Note: Any nonzero constant term of `V` is ignored.
		Z.Append(VV[0].Times(U.Subst(VV[1])))

	}(Z, U, V)
	return Z
}

// MonSubst Monomial Substition: `U(c*x^n)`
// Each Ui is multiplied by `c^i` and followed by n-1 zeros.
// If `n` is not positive, zero (an empty closed power series) is returned.
func (U PS) MonSubst(c0 Coefficient, n int) PS {
	Z := U.new()

	if !(n > 0) {
		Z.Close()
		U.Drop()
		return Z
	}

	go func(Z PS, U PS, c0 Coefficient, n int) {
		c := aOne()
		for Z.SendCfnFrom(U, cMul(c)) { // `c * u`
			for i := 1; i < n; i++ {
				if !Z.Put(aZero()) { // n-1 zeros
					U.Drop()
					return
				}
			}
			c.Mul(c, c0) // `c = c * c0 = c^i`
		}
	}(Z, U, c0, n)
	return Z
}

// ===========================================================================
