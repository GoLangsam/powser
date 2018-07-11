// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package ps

import (
	"github.com/GoLangsam/powser/big"
)

// GetValS returns an array with each first value received from the given power series.
//
// BUG: As of now, it works for pairs only!
func (in PS2) GetValS() (out [2]*big.Rat) {
	n := len(in)
	if n != 2 {
		panic("bad n in Get2")
	}

	req := make([]chan<- struct{}, 0, n)
	snd := make([]<-chan *big.Rat, 0, n) // we might send here
	dat := make([]<-chan *big.Rat, 0, n) // we shall send here

	for i := 0; i < n; i++ {
		req[i], snd[i] = in[i].From()
		dat[i] = nil
	}

	for n = 2 * n; n > 0; n-- {

		select {
		// whish we could repeat this n times
		case req[0] <- struct{}{}:
			dat[0] = snd[0]
			req[0] = nil
		case req[1] <- struct{}{}:
			dat[1] = snd[1]
			req[1] = nil
		// whish we could repeat this n times
		case it := <-dat[0]:
			out[0] = it
			dat[0] = nil
		case it := <-dat[1]:
			out[1] = it
			dat[1] = nil
		}
	}
	return out
}
