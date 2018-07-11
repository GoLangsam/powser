// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// Split returns a pair of power series identical to a given power series.
func (U PS) Split() PS2 {
	UU := U.NewPair()
	UU.Split(U)
	return UU
}

// Split inp into a given pair of power series.
func (out PS2) Split(in PS) {
	release := make(chan struct{})
	go split(out[0], out[1], in, release)
	close(release)
}

// split reads a single demand channel and replicates its
// output onto two, which may be read at different rates.
//
// A process is created at first demand and dies
// after the data has been sent to both outputs.
//
// When multiple generations of split exist,
// the newest will service requests on one channel,
// which is always renamed to be out[0];
// the elder will service requests on the other channel, out[1].
// All generations but the newest hold queued data
// that has already been sent to out[0].
//
// When data has finally been sent to out[1] by the newest generation,
// a signal on the release-wait channel tells the next newer
// generation to begin servicing out[1].
//
func split(out1, out2, inp PS, wait <-chan struct{}) {

	req1, _ := out1.Into()
	req2, _ := out2.Into()

	both := false // do not service both channels before <-wait
	req := false  // got valid request?

	select {
	case _, req = <-req1:

	case <-wait:
		both = true
		select {
		case _, req = <-req1:

		case _, req = <-req2: // swap
			out1, out2 = out2, out1
			// req1, req2 = req2, req1
		}
	}

	dat, ok := inp.Get()
	if !ok {
		inp.Drop()
	}

	release := make(chan struct{})

	if ok && req { // dispatch - as we have data
		go split(out1, out2, inp, release)
		out1.Send(dat)
	} else {
		out1.Close()
	}

	if !both {
		<-wait
	}

	if next := out2.Next(); next && ok {
		out2.Send(dat)
	} else if !(ok && req) { // no dispatch - we're last man standing
		out2.Close()
	}

	close(release)
}
