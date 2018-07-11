// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// ===========================================================================

// get2 returns each first value received from the two given power series
// together with its respective ok boolean.
func get2(inp1, inp2 PS) (Coefficient, bool, Coefficient, bool) {

	req1, snd1 := inp1.From()
	req2, snd2 := inp2.From()
	dat1, dat2 := snd1, snd2

	dat1, dat2 = nil, nil

	var out1, out2 Coefficient
	var oks1, oks2 bool

	for i := 2 * 2; i > 0; i-- {

		select {

		case req1 <- struct{}{}:
			dat1, req1 = snd1, nil
		case req2 <- struct{}{}:
			dat2, req2 = snd2, nil

		case out1, oks1 = <-dat1:
			dat1 = nil
		case out2, oks2 = <-dat2:
			dat2 = nil

		}
	}
	return out1, oks1, out2, oks2
}

// getValS returns a slice with each first value received from the given power series.
//
// BUG: As of now, it works for pairs only!
func getValS(in ...PS) ([]Coefficient, []bool) {
	n := len(in)
	if n != 2 {
		panic("getValS must have exactly 2 arguments")
	}

	req := make([]chan<- struct{}, 0, n)    // we request here - initially
	snd := make([]<-chan Coefficient, 0, n) // we might receive here
	dat := make([]<-chan Coefficient, 0, n) // we shall receive here
	oks := make([]bool, 0, n)               // did we receive here?
	out := make([]Coefficient, 0, n)        // the values to be returned

	for i := 0; i < n; i++ {
		req[i], snd[i] = in[i].From() // from
		dat[i] = nil                  // block receive on dat initially
		out[i] = NewCoefficient(0, 0) // init
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
