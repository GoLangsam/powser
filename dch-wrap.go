// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================
// Wrappers for multi-argument methods
// use dch.MyDchInto() and dch.MyDchFrom() to obtain the anonymously embedded value
// and invoke its underlying method.

// Append all coefficients from `From` into `Into`.
func (Into PS) Append(From PS) {
	Into.MyDch().Append(From.MyDch())
}

// GetNextFrom `From` for `Into` and report success.
// Follow with `Into.Send( f(c) )`, iff ok.
func (Into PS) GetNextFrom(From PS) (c Coefficient, ok bool) {
	return Into.MyDch().GetNextFrom(From.MyDch())
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

// Split returns a pair of power series identical to the given one.
func (U PS) Split() [2]PS {
	UU := U.newPair()
	U.MyDch().SplitUs(UU[0].MyDch(), UU[1].MyDch())
	return UU
}

// ---------------------------------------------------------------------------

// pairPS represents a pair of power series.
type pairPS [2]PS

// pair returns an empty pair of new power series.
func (U PS) newPair() pairPS {
	return pairPS{New(), New()}
}

// Split `From` into a given pair of power series.
func (UU pairPS) Split(From PS) {
	From.MyDch().SplitUs(UU[0].MyDch(), UU[1].MyDch())
}

// ===========================================================================
