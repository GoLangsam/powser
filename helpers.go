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
// Note: Its used where it helps to tighten the implementation of an algorithm,
// and intentionally calculations are done directly and explicit in other places.

// cSame `u`
func cSame() func(u Coefficient) Coefficient {
	return func(u Coefficient) Coefficient {
		return u
	}
}

// cRatIby1 `u * i`
func cRatIby1(i int) func(u Coefficient) Coefficient {
	return func(u Coefficient) Coefficient {
		return aC().Mul(u, ratIby1(i))
	}
}

// cRat1byI `u * 1/i`
func cRat1byI(i int) func(u Coefficient) Coefficient {
	return func(u Coefficient) Coefficient {
		return aC().Mul(u, rat1byI(i))
	}
}

// cMul `c * u`
func cMul(c Coefficient) func(u Coefficient) Coefficient {
	return func(u Coefficient) Coefficient {
		return aC().Mul(c, u)
	}
}

// SendCfnFrom `cfn(From)` into `Into` and report success.
func (Into PS) SendCfnFrom(From PS, cfn func(c Coefficient) Coefficient) (ok bool) {
	var c Coefficient
	if c, ok = Into.GetNextFrom(From); ok {
		Into.Send(cfn(c))
	}
	return
}

// ===========================================================================
