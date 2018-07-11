// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ps

// split reads a single demand channel and replicates its
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
func split(out1, out2, inp PS, wait <-chan struct{}) {
	both := false // do not service both channels

	req1, _ := out1.Into()
	req2, _ := out2.Into()

	var req bool // got request?

	select {
	case _, req = <-req1:

	case <-wait:
		both = true
		select {
		case _, req = <-req1:

		case _, req = <-req2: // swap
			out1, out2 = out2, out1
			req1, req2 = req2, req1
		}
	}

	if dat, ok := inp.Get(); ok {
		release := make(chan struct{})
		go split(out1, out2, inp, release)

		if req {
			out1.Send(dat)
		}
		if !both {
			<-wait
		}
		if out2.Next() {
			out2.Send(dat)
		}

		close(release)

	} else {
		out1.Close()
		out2.Close()
	}
}
