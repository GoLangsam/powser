// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dch

import (
	"math/big"
	// "github.com/GoLangsam/powser/big"
)

// ===========================================================================

// GetWith returns each first value received from the two given demand channels
// together with their respective ok boolean.
func (from *Dch) GetWith(with *Dch) (valFrom *big.Rat, okFrom bool, valWith *big.Rat, okWith bool) {

	reqFrom, sndFrom := from.req, from.ch
	reqWith, sndWith := with.req, with.ch
	datFrom, datWith := sndFrom, sndWith
	datFrom, datWith = nil, nil // block receives initially

	for i := 2 * 2; i > 0; i-- {

		select {

		case reqFrom <- struct{}{}:
			datFrom, reqFrom = sndFrom, nil
		case reqWith <- struct{}{}:
			datWith, reqWith = sndWith, nil

		case valFrom, okFrom = <-datFrom:
			datFrom = nil
		case valWith, okWith = <-datWith:
			datWith = nil

		}
	}
	return
}

// getValS returns a slice with each first value received from the given demand channels.
//
// BUG: As of now, it works for pairs only!
func getValS(in ...*Dch) ([]*big.Rat, []bool) {
	n := len(in)
	if n != 2 {
		panic("getValS must have exactly 2 arguments")
	}

	req := make([]chan<- struct{}, n) // we request here - initially
	snd := make([]<-chan *big.Rat, n) // we might receive here
	dat := make([]<-chan *big.Rat, n) // we shall receive here
	oks := make([]bool, n)            // did we receive here?
	out := make([]*big.Rat, n)        // the values to be returned

	for i := 0; i < n; i++ {
		req[i], snd[i] = in[i].From() // from
		dat[i] = nil                  // block receive on dat initially
	}

	for n = 2 * n; n > 0; n-- {

		select {

		// whish we could repeat this n times
		case req[0] <- struct{}{}:
			dat[0], req[0] = snd[0], nil // open dat & block req
		case req[1] <- struct{}{}:
			dat[1], req[1] = snd[1], nil // open dat & block req

		// whish we could repeat this n times
		case out[0], oks[0] = <-dat[0]:
			dat[0] = nil // block receive on dat again
		case out[1], oks[1] = <-dat[1]:
			dat[1] = nil // block receive on dat again

		}
	}
	return out, oks
}

// ===========================================================================
