// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================
// Wrappers for multi-argument methods
// use dch.Self() to obtain the anonymously embedded value
// and invoke its underlying method.

// Append all coefficients from `U` into `Into`.
func (Into PS) Append(U PS) {
	Into.Self().Append(U.Self())
}

// NextGetFrom `U` for `Into` and report success.
// Follow with `Into.Send( f(c) )`, iff ok.
func (Into PS) NextGetFrom(U PS) (c Coefficient, ok bool) {
	return Into.Self().NextGetFrom(U.Self())
}

// GetWith returns each first value received from the two given power series
// together with their respective ok boolean.
func (U PS) GetWith(V PS) (cU Coefficient, okU bool, cV Coefficient, okV bool) {
	cU, okU, cV, okV = U.Self().GetWith(V.Self())
	if !okU {
		cU = aZero()
	}
	if !okV {
		cV = aZero()
	}

	return
}

// ---------------------------------------------------------------------------

// Split returns two power series identical to (the current remainder of) the given `from` power series.
// It's a convenient wrapper around SplitUs.
func (U PS) Split() (PS, PS) {
	S1, S2 := U.Self().Split()
	into1, into2 := PS{S1}, PS{S2}
	return into1, into2
}

// SplitUs from `from` into two given power series.
func (U PS) SplitUs(V, W PS) {
	U.Self().SplitUs(V.Self(), W.Self())
}

// ===========================================================================
