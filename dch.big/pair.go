// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dch

// Pairing reads a single demand channel and replicates its
// output onto two, which may be read at different rates
// (no lock-step).
//
// A process is created at first demand and dies
// after the data has been sent to both outputs.
//
// When multiple generations of split exist,
// the newest will service requests on one channel,
// which is always renamed to be `out1`.
// All elder wait to service requests
// on the other channel, `out2`.
//
// All generations but the newest hold queued data
// that has already been sent to `out1`.
//
// When data has finally been sent to `out2` also,
// a signal on the release-wait channel tells
// the next newer (one step elder) generation
// to begin servicing `out2`.
//

func (from *Dch) Pair() (out1, out2 *Dch) {
	cha1 := New()
	cha2 := New()
	go from.Split(cha1, cha2)
	return cha1, cha2
}

func (from *Dch) Split(out1, out2 *Dch) {

	release := make(chan struct{})
	go from.split(out1, out2, release)
	release <- struct{}{}
}

func (from *Dch) split(out1, out2 *Dch, wait <-chan struct{}) {
	both := false // do not service both channels

	reqI, datI := from.From()
	req0, dat0 := out1.Into()
	req2, dat1 := out2.Into()

	select {
	case <-req0:

	case <-wait:
		both = true
		select {
		case <-req0:

		case <-req2: // swap
			out1, out2 = out2, out1
			req0, req2 = req2, req0
			dat0, dat1 = dat1, dat0
		}
	}

	reqI <- struct{}{}
	release := make(chan struct{})
	go from.split(out1, out2, release)
	dat := <-datI
	dat0 <- dat
	if !both {
		<-wait
	}
	<-req2
	dat1 <- dat
	release <- struct{}{}
}
