// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================
// Methods and functions helpful for power series and pairs thereof.

// Conventions
// Upper-case for power series.
// Lower-case for coefficients.
// Input variables:  From U,V,...
// Output variables: Into ...,Y,Z

// new returns a fresh power series.
func (U PS) new() PS {
	return New()
}

// ---------------------------------------------------------------------------
// Closure functions on some coefficient math - for convenient use with SendCfnFrom
// Note: Such closures are used where it helps to tighten the implementation of an algorithm,
// and in other places calculations are intentionally done directly and explicit.

// CoefficientFunc represents a function which -given some Coefficient- returns some new(!) related Coefficient.
// Functions which return a CoefficientFunc (as a closure)
//  - are used with SendCfnFrom to simplify some arithmetic
//  - help to avoid pitfalls related to accidental re-use of method receivers
//    (a common 'challenge' when using *big.Rat)
type CoefficientFunc func(Coefficient) Coefficient

// cSame `u`
func cSame() CoefficientFunc {
	return func(u Coefficient) Coefficient {
		return u
	}
}

// cAdd `c + u`
func cAdd(c Coefficient) CoefficientFunc {
	return func(u Coefficient) Coefficient {
		return aC().Add(c, u)
	}
}

// cNeg `-u`
func cNeg() CoefficientFunc {
	return func(u Coefficient) Coefficient {
		return aC().Neg(u)
	}
}

// cRatIby1 `u * i`
func cRatIby1(i int) CoefficientFunc {
	return func(u Coefficient) Coefficient {
		return aC().Mul(u, ratIby1(i))
	}
}

// cRat1byI `u * 1/i`
func cRat1byI(i int) CoefficientFunc {
	return func(u Coefficient) Coefficient {
		return aC().Mul(u, rat1byI(i))
	}
}

// cMul `c * u`
func cMul(c Coefficient) CoefficientFunc {
	return func(u Coefficient) Coefficient {
		return aC().Mul(c, u)
	}
}

// cInv `1/ u`
func cInv() func(u Coefficient) Coefficient {
	return func(u Coefficient) Coefficient {
		return aC().Inv(u)
	}
}

// SendCfnFrom sends into `Into` the result of the CoefficientFunc `cfn`
// applied to what can be `Get` from `From` -`cfn(From)` so to say-
// and report success.
//
// SendCfnFrom includes housekeeping (implicitly via GetNextFrom):
// If `into` has been dropped, `from` is dropped and `into` is closed,
// and failure of normal progress is reported: ok == false.
func (Into PS) SendCfnFrom(From PS, cfn CoefficientFunc) (ok bool) {
	var c Coefficient
	if c, ok = Into.NextGetFrom(From); ok {
		ok = Into.Send(cfn(c))
	}
	return
}

// ===========================================================================
