// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package ps

// ===========================================================================
// Specific power series

/*
--- https://en.wikipedia.org/wiki/Formal_power_series

1 / (1-x)	<=> a(n) = 1
1 / (1+x)	<=> a(n) = (-1)^n
x / (1-x)^2	<=> a(n) = n

--- e^x
U := a0 a1 a2 ...
e^x = c0 c1 c2 ...
- ci = 1 / i!

--- Sinus
U := a0 a1 a2 ...
Sin(U) = c0 c1 c2 ...
- c2i = 0
- c2i+1 = (-1)^i / (2i+1)!

--- Cosinus
U := a0 a1 a2 ...
Cos(U) = c0 c1 c2 ...
- c2i = (-1)^i / (2i)!
- c2i+1 = 0

--- Power-Exponentiation
U := a0 a1 a2 ...
U^n = c0 c1 c2 ...
- c0 = a0^n
- ci = 1/(i*a0) * sum k=1...i[ ( k*i - i + k ) * ak * ci-k ]

--- Division
U := a0 a1 a2 ...
V := b0 b1 b2 ...
U/V = c0 c1 c2 ...
- c0 = 1 / a0
- ci = (1/b0) * ( ai - sum k=1...i[ bk - ci-k ] )

--- Exponential
U := a0 a1 a2 ...
Exp(U) = c0 c1 c2 ...
- c0 = 1
- ci = 1 / i!

--- Composition
Subst == Composition

--- https://en.wikipedia.org/wiki/Formal_power_series#The_Lagrange_inversion_formula

*/

// ===========================================================================

// Ones are 1 1 1 1 1 ... = `1/(1-x)` with a simple pole at `x=1`.
func Ones() PS {
	return AdInfinitum(NewCoefficient(1, 1))
}

// Twos are 2 2 2 2 2 ... just for samples.
func Twos() PS {
	return AdInfinitum(NewCoefficient(2, 1))
}

// AdInfinitum repeats coefficient `c` ad infinitum
// and returns `c^i`.
func AdInfinitum(c Coefficient) PS {
	Z := New()
	go func(Z PS, c Coefficient) {
		for Z.Put(c) {
		}
	}(Z, c)
	return Z
}

// ===========================================================================

// Factorials starting from zero: 1, 1, 2, 6, 24, 120, 720, 5040 ...
func Factorials() PS {
	Z := New()
	go func(Z PS) {
		curr := aOne()
		for i := 1; Z.Put(curr); i++ {
			curr = curr.Mul(curr, ratIby1(i))
		}
	}(Z)
	return Z
}

// OneByFactorial starting from zero: 1/1, 1/1, 1/2, 1/6, 1/120 ...
func OneByFactorial() PS {
	Z := New()
	go func(Z PS) {
		curr := aOne()
		for i := 1; Z.Put(aC().Inv(curr)); i++ {
			curr = curr.Mul(curr, ratIby1(i))
		}
	}(Z)
	return Z
}

// Fibonaccis starting from zero: 1, 2, 3, 5, 8, 13, 21, 34, 55, 89 ...
func Fibonaccis() PS {
	Z := New()
	go func(Z PS) {
		prev, curr := aZero(), aOne()
		for Z.Put(curr) {
			prev, curr = curr, aC().Add(curr, prev)
		}
	}(Z)
	return Z
}

// OneByFibonacci starting from zero: 1/1, 1/2, 1/3, 1/5, 1/8, 1/13 ...
func OneByFibonacci() PS {
	Z := New()
	go func(Z PS) {
		prev, curr := aZero(), aOne()
		for Z.Put(aC().Inv(curr)) {
			prev, curr = curr, aC().Add(curr, prev)
		}
	}(Z)
	return Z
}

// Harmonics: 1, 1+ 1/2, 1+ 1/2+ 1/3, 1+ 1/2+ 1/3+ 1/4 ...
//  `1/(1-x) * ln( 1/(1-x) )`
func Harmonics() PS {
	Z := New()
	go func(Z PS) {
		curr := aOne()
		for i := 2; Z.Put(curr); i++ {
			curr = aC().Add(curr, rat1byI(i))
		}
	}(Z)
	return Z
}

// Sincos returns the power series for sine and cosine (in radians).
func Sincos() (Sin PS, Cos PS) {
	Sin = New()
	Cos = New()

	U := OneByFactorial()
	UU := U.Split()

	f := func(Z PS, U PS, odd bool) {
		var minus bool
		for {
			if u, ok := Z.NextGetFrom(U); ok {
				if odd {
					if minus {
						Z.Send(u.Neg(u))
					} else {
						Z.Send(u)
					}
					minus = !minus
				} else {
					Z.Send(aZero())
				}
				odd = !odd
			} else {
				return
			}
		}
	}

	go f(Sin, UU[0], false)
	go f(Cos, UU[1], true)

	return
}

// Sin returns the power series for sine (in radians).
func Sin() PS {
	U, V := Sincos()
	V.Drop()
	return U
}

// Cos returns the power series for cosine (in radians).
func Cos() PS {
	U, V := Sincos()
	U.Drop()
	return V
}

// ===========================================================================
