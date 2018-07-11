// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================

// GetNextFrom `From` for `Into` and report success.
// Follow with `Into.Send( f(c) )`, iff ok.
func (Into PS) GetNextFrom(From PS) (c Coefficient, ok bool) {
	return Into.MyDch().GetNextFrom(From.MyDch())
}

// Append all coefficients from `From` into `Into`.
func (Into PS) Append(From PS) {
	Into.MyDch().Append(From.MyDch())
}

// append all coefficients from `From` into `Into`.
// without cleanup of handshaking resources.
func (Into PS) append(From PS) {
	Into.MyDch().AppendOnly(From.MyDch())
}

// Split returns a pair of power series identical to (the current remainder of) a given power series.
func (U PS) Split() [2]PS {
	UU := U.pair()
	U.MyDch().SplitUs(UU[0].MyDch(), UU[1].MyDch())
	return UU
}

// Split `inp` into a given pair of power series.
func (out pairPS) Split(inp PS) {
	inp.MyDch().SplitUs(out[0].MyDch(), out[1].MyDch())
}

// GetWith returns each first value received from the two given power series
// together with their respective ok boolean.
func (U PS) GetWith(V PS) (cU Coefficient, okU bool, cV Coefficient, okV bool) {
	cU, okU, cV, okV = U.MyDch().GetWith(V.MyDch())
	if !okU {
		cU = aZero()
	}
	if !okV {
		cV = aZero()
	}

	return
}

// ===========================================================================
