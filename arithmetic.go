// Copyright 2017 Andreas Pannewitz. All rights reserved.
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

// New returns a fresh power series.
func (U PS) New() PS {
	return NewPS()
}

// NewPair returns an empty pair of new power series.
func (U PS) NewPair() PS2 {
	return PS2{NewPS(), NewPS()}
}

// Obtain for `Into` from `From` and report success.
func (Into PS) Obtain(From PS) (c Coefficient, ok bool) {
	if ok = Into.Next(); ok {
		c, ok = From.Get()
	}
	if !ok {
		From.Drop()
		Into.Close()
	}
	return
}

// SendOneFrom for `Into` from `From` and report success.
func (Into PS) SendOneFrom(From PS) (ok bool) {
	var c Coefficient
	if c, ok = Into.Obtain(From); ok {
		Into.Send(c)
	}
	return
}

// Clone returns a new powerseries to receive values from `U`.
func (From PS) Clone() PS {
	Into := From.New()
	go Into.Append(From)
	return Into
}

// Append all coefficients from `Z` to `U`.
func (Into PS) Append(From PS) {
	defer func() {
		From.Drop()
		Into.Close()
	}()
	Into.append(From)
}

// append all coefficients from `Z` to `U`
// without cleanup of handshaking resources.
func (Into PS) append(From PS) {
	var c Coefficient
	var ok bool
	for Into.Next() {
		if c, ok = From.Get(); !ok {
			return
		}
		Into.Send(c)
	}
}

// ===========================================================================
// Add two power series.
func Add(U, V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {

		var u, v Coefficient
		var uok, vok bool
		for Z.Next() {
			u, uok, v, vok = get2(U, V)
			switch { // fini(u) + 2*fini(v) {
			case uok && vok:
				Z.Send(u.Add(u, v))
			case uok:
				V.Drop()
				Z.Send(u)
				Z.Append(U)
			case vok:
				U.Drop()
				Z.Send(v)
				Z.Append(V)
			default:
				U.Drop()
				V.Drop()
				Z.Close()
				return
			}
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
	return U.Plus(V.Cmul(aMinusOne))
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

// Cmul multiplies `U` by a constant `c`
// and returns `c*U`.
func (U PS) Cmul(c Coefficient) PS {
	Z := U.New()
	go func(c Coefficient, U, Z PS) {
		for u, ok := Z.Obtain(U); ok; u, ok = Z.Obtain(U) {
			Z.Send(u.Mul(c, u))
		}
	}(c, U, Z)
	return Z
}

// Monmul multiplies `U` by the monomial "x^n"
// and returns `x^n * U`.
func (U PS) Monmul(n int) PS {
	Z := U.New()
	go func(n int, U PS, Z PS) {
		for ; n > 0; n-- {
			Z.Put(aZero)
		}
		Z.Append(U)
	}(n, U, Z)
	return Z
}

// Xmul multiplies `U` by `x`
// (by the monomial "x^1")
// and returns `x * U`.
func (U PS) Xmul() PS {
	return U.Monmul(1)
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
	Z := NewPS()
	go func(U, V, Z PS) {

		if !Z.Next() {
			U.Drop()
			V.Drop()
			Z.Close()
			return
		}
		u, uok, v, vok := get2(U, V)
		if !uok || !vok {
			U.Drop()
			V.Drop()
			Z.Close()
			return
		}

		c := u // `u*v`
		c.Mul(u, v)
		Z.Send(c)

		UU := U.Split()
		VV := V.Split()

		W1 := Add(VV[0].Cmul(u), UU[0].Cmul(v))

		if Z.SendOneFrom(W1) {
			Z.Append(W1.Plus(Mul(UU[1], VV[1])))
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
		u, ok := Z.Obtain(U)
		if !ok {
			return
		}
		// constant term: drop
		// Thus: we must Z.Send() before another Z.Next()
		// and may not use an Obtain-loop and have to cleanup ourselfs

		i := 1
		for u, ok = U.Get(); ok; u, ok = U.Get() {
			Z.Send(u.Mul(ratIby1(i), u))
			if !Z.Next() {
				Z.Close()
				U.Drop()
				return
			}
			i++
		}

	}(U, Z)
	return Z
}

// Integrate, with const of integration.
func (U PS) Integ(c Coefficient) PS {
	Z := U.New()
	go func(c Coefficient, U, Z PS) {
		Z.Put(c)

		i := 1
		for u, ok := Z.Obtain(U); ok; u, ok = Z.Obtain(U) {
			Z.Send(u.Mul(rat1byI(i), u))
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
		z, ok := Z.Obtain(U)
		if !ok {
			return
		}
		Z.Send(z.Inv(z)) // `1/u`

		mz := z
		mz.Neg(z) // minus z `-z`
		ZZ := U.NewPair()
		ZZ.Split(Mul(U.Cmul(mz), ZZ[0].Shift(z)))
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
	ZZ := U.NewPair()
	ZZ.Split(Mul(ZZ[0], U.Deriv()).Integ(aOne))
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
		u, ok := Z.Obtain(U)
		if !ok {
			V.Drop()
			return
		}
		Z.Send(u)

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
		c := aOne
		uc := c // `u * c`
		for u, ok := Z.Obtain(U); ok; u, ok = Z.Obtain(U) {
			Z.Send(uc.Mul(u, c))
			c.Mul(c, c0)

			for i := 1; i < n; i++ {
				if !Z.Next() {
					Z.Close()
					U.Drop()
					return
				}
				Z.Send(aZero)
			}
		}
	}(U, Z, c0, n)
	return Z
}

// ===========================================================================
