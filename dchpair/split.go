// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package dchpair

import (
	. "github.com/GoLangsam/powser/dch"
)

// Split reads a single demand channel and replicates its
// output onto two, which may be read at different rates.
// A process is created at first demand and dies
// after the data has been sent to both outputs.
//
// When multiple generations of split exist, the newest
// will service requests on one channel, which is
// always renamed to be out[0]; the oldest will service
// requests on the other channel, out[1].  All generations but the
// newest hold queued data that has already been sent to
// out[0].  When data has finally been sent to out[1],
// a signal on the release-wait channel tells the next newer
// generation to begin servicing out[1].
//
func (out DchPair) Split(in *Dch) {
	release := make(chan struct{})
	go out.dosplit(in, release)
	release <- struct{}{}
}

func (out DchPair) dosplit(in *Dch, wait <-chan struct{}) {
	both := false // do not service both channels

	reqI, datI := in.From()
	req0, dat0 := out[0].Into()
	req1, dat1 := out[1].Into()

	select {
	case <-req0:

	case <-wait:
		both = true
		select {
		case <-req0:

		case <-req1: // swap
			out[0], out[1] = out[1], out[0]
			req0, req1 = req1, req0
			dat0, dat1 = dat1, dat0
		}
	}

	reqI <- struct{}{}
	release := make(chan struct{})
	go out.dosplit(in, release)
	dat := <-datI
	dat0 <- dat
	if !both {
		<-wait
	}
	<-req1
	dat1 <- dat
	release <- struct{}{}
}
