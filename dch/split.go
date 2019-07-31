// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dch

// ===========================================================================

// Split returns two demand channels identical to (the current remainder of) the given `from` channel.
// It's a convenient wrapper around SplitUs.
func (from *Dch) Split() (*Dch, *Dch) {
	into1, into2 := from.New(), from.New()
	from.SplitUs(into1, into2)
	return into1, into2
}

// SplitUs from `from` into two given demand channels.
func (from *Dch) SplitUs(into1, into2 *Dch) {
	release := make(chan struct{})
	go from.split(into1, into2, release)
	close(release)
}

// ===========================================================================

// split reads a single demand channel and replicates its
// output onto two, which may be read at different rates.
//
// A process is created at first demand and dies
// after the data has been sent to both outputs.
//
// When multiple generations of split exist,
// the newest will service requests on one channel,
// (which is always renamed to be `out1`);
// the elder will service requests on the other channel, (`out2`).
//
// All generations but the newest hold queued data
// that has already been sent to `out1`.
//
// When data has finally been sent to `out2` by the newest generation,
// a signal on the release-wait channel tells the next newer
// generation to begin servicing `out2`.
//
// When `inp` becomes closed or `out1` ceases to request,
// `out1` will be closed, and -after the queue has emptied-
// the last living process will append `inp` to `out2`.
func (from *Dch) split(out1, out2 *Dch, wait <-chan struct{}) {

	req1 := out1.req
	req2 := out2.req

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

	dat, ok := from.Get()

	release := make(chan struct{})
	if ok && req { // dispatch - as we have data
		go from.split(out1, out2, release)
		out1.Send(dat)
	} else {
		out1.Close()
	}

	if !both {
		<-wait
	}

	if ok {
		out2.Provide(dat)
	}
	if !(ok && req) { // no dispatch - we're last man standing
		out2.Append(from)
	}

	close(release)
}

// ===========================================================================
