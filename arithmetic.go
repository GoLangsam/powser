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
		var u, v Coefficient
		for {
			<-Z_req
			u, v = GetVal2(U, V)
			switch fini(u) + 2*fini(v) {
			case 0:
				Z_in <- u.Add(u, u)
			case 1:
				Z_in <- v
				Z.Append(V)
			case 2:
				Z_in <- u
				Z.Append(U)
			case 3:
				Z_in <- finis
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
		done := false
		for !done {
			<-Z_req
			u := U.Get()
			if atEnd(u) {
				done = true
			} else {
				Z_in <- u.Mul(c, u)
			}
		}
		Z_in <- finis
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
		<-Z_req
		u, v := GetVal2(U, V)
		if atEnd(u) || atEnd(v) {
			Z_in <- finis
		} else {
			var prod Coefficient
			prod.Mul(u, v)
			Z_in <- prod
			UU := U.Split()
			VV := V.Split()
			W := Add(VV[0].Cmul(u), UU[0].Cmul(v))
			<-Z_req
			Z_in <- W.Get()
			Z.Append(W.Plus(Mul(UU[1], VV[1])))
		}
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
		Z_req, Z_dat := Z.Into()
		<-Z_req
		u := U.Get()
		if !atEnd(u) {
			done := false
			for i := 1; !done; i++ {
				u = U.Get()
				if atEnd(u) {
					done = true
				} else {
					Z_dat <- u.Mul(ratI(i), u)
					<-Z_req
				}
			}
		}
		Z_dat <- finis
	}(U, Z)
	return Z
}

// Integrate, with const of integration
func (U PS) Integ(c Coefficient) PS {
	Z := NewPS()
	go func(c Coefficient, U, Z PS) {
		Z_req, Z_in := Z.Into()
		Z.Put(c)
		done := false
		for i := 1; !done; i++ {
			<-Z_req
			u := U.Get()
			if !atEnd(u) {
				Z_in <- u.Mul(rat1byI(i), u)
			} else {
				Z_in <- finis
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
//	ZZ = `1/u * -UU * (z + x*ZZ)`
func (U PS) Recip() PS {
	Z := NewPS()
	go func(U, Z PS) {
		Z_req, Z_in := Z.Into()
		ZZ := NewPS2()
		<-Z_req
		z := U.Get()
		z.Inv(z)
		Z_in <- z
		var mz Coefficient
		mz.Neg(z) // minus z `-z`
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
	ZZ := NewPS2()
	ZZ.Split(Mul(ZZ[0], U.Diff()).Integ(aOne))
	return ZZ[1]
}

// Subst itute V for x in U, where the leading term of V is zero
//	let U = `u + x*UU`
//	let V = `v + x*VV`
//	then Subst(U,V) = `u + VV*S(V,UU)`
// Note: Any nonzero constant term of `V` is ignored.
func (U PS) Subst(V PS) PS {
	Z := NewPS()
	go func(U, V, Z PS) {
		Z_req, Z_in := Z.Into()
		VV := V.Split()
		<-Z_req
		u := U.Get()
		Z_in <- u
		if !atEnd(u) {
			if !atEnd(VV[0].Get()) {
				Z.Append(Mul(VV[0], U.Subst(VV[1])))
			} else {
				Z.Put(finis)
			}
		}
	}(U, V, Z)
	return Z
}

// MonSubst Monomial Substition: U(c x^n)
// Each Ui is multiplied by `c^i` and followed by n-1 zeros.
func (U PS) MonSubst(c0 Coefficient, n int) PS {
	Z := NewPS()
	go func(U, Z PS, c0 Coefficient, n int) {
		Z_req, Z_in := Z.Into()
		c := aOne
		var uc Coefficient // `u * c`
		for {
			<-Z_req
			u := U.Get()
			Z_in <- uc.Mul(u, c)
			c.Mul(c, c0)
			if atEnd(u) {
				Z_in <- finis
				break
			}
			for i := 1; i < n; i++ {
				<-Z_req
				Z_in <- aZero
			}
		}
	}(U, Z, c0, n)
	return Z
}

// ===========================================================================
