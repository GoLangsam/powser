// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package dch

import (
	"github.com/GoLangsam/powser/rat"
)

type DchPair [2]*Dch

func NewPair() (pair DchPair) {
	pair[0] = New()
	pair[1] = New()
	return pair
}

// split reads a single demand channel and replicates its
// output onto two, which may be read at different rates.
// A process is created at first demand for a rat and dies
// after the rat has been sent to both outputs.

// When multiple generations of split exist, the newest
// will service requests on one channel, which is
// always renamed to be out[0]; the oldest will service
// requests on the other channel, out[1].  All generations but the
// newest hold queued data that has already been sent to
// out[0].  When data has finally been sent to out[1],
// a signal on the release-wait channel tells the next newer
// generation to begin servicing out[1].

func (out DchPair) Split(in *Dch) {
	release := make(chan struct{})
	go out.dosplit(in, release)
	release <- struct{}{}
}

func (out DchPair) dosplit(in *Dch, wait <-chan struct{}) {
	both := false // do not service both channels

	select {
	case <-out[0].req:

	case <-wait:
		both = true
		select {
		case <-out[0].req:

		case <-out[1].req:
			out[0], out[1] = out[1], out[0]
		}
	}

	seqno++
	in.req <- seqno
	release := make(chan struct{})
	go out.dosplit(in, release)
	dat := <-in.dat
	out[0].dat <- dat
	if !both {
		<-wait
	}
	<-out[1].req
	out[1].dat <- dat
	release <- struct{}{}
}

// Get one rat from each of 2 demand channels
func Get2(in DchPair) (out [2]*rat.Rat) {
	n := len(in)
	if n != 2 {
		panic("bad n in getn")
	}
	req := new([2]chan int)
	dat := new([2]chan *rat.Rat)

	var i int
	var it *rat.Rat
	for i = 0; i < n; i++ {
		req[i] = in[i].req
		dat[i] = nil
	}
	for n = 2 * n; n > 0; n-- {
		seqno++

		select {
		case req[0] <- seqno:
			dat[0] = in[0].dat
			req[0] = nil
		case req[1] <- seqno:
			dat[1] = in[1].dat
			req[1] = nil
		case it = <-dat[0]:
			out[0] = it
			dat[0] = nil
		case it = <-dat[1]:
			out[1] = it
			dat[1] = nil
		}
	}
	return out
}
