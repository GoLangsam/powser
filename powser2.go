// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

// Like powser1.go but uses channels of interfaces.
// Has not been cleaned up as much as powser1.go, to keep
// it distinct and therefore a different test.

// Power series package
// A power series is a channel, along which flow rational
// coefficients.  A denominator of zero signifies the end.
// Original code in Newsqueak by Doug McIlroy.
// See Squinting at Power Series by Doug McIlroy,
//   http://www.cs.bell-labs.com/who/rsc/thread/squint.pdf

package ps

import (
	"github.com/GoLangsam/powser/rat"
)

type dch struct {
	req chan int
	dat chan *rat.Rat
	nam int
}

type dch2 [2]*dch

var chnames string
var chnameserial int
var seqno int

func mkdch() *dch {
	c := chnameserial % len(chnames)
	chnameserial++
	d := new(dch)
	d.req = make(chan int)
	d.dat = make(chan *rat.Rat)
	d.nam = c
	return d
}

func mkdch2() *dch2 {
	d2 := new(dch2)
	d2[0] = mkdch()
	d2[1] = mkdch()
	return d2
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

func dosplit(in *dch, out *dch2, wait chan int) {
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
	release := make(chan int)
	go dosplit(in, out, release)
	dat := <-in.dat
	out[0].dat <- dat
	if !both {
		<-wait
	}
	<-out[1].req
	out[1].dat <- dat
	release <- 0
}

func split(in *dch, out *dch2) {
	release := make(chan int)
	go dosplit(in, out, release)
	release <- 0
}

func put(dat *rat.Rat, out *dch) {
	<-out.req
	out.dat <- dat
}

func get(in *dch) *rat.Rat {
	seqno++
	in.req <- seqno
	return <-in.dat
}

// Get one item from each of n demand channels

func getn(in [2]*dch) [2]*rat.Rat {
	n := len(in)
	if n != 2 {
		panic("bad n in getn")
	}
	req := new([2]chan int)
	dat := new([2]chan *rat.Rat)
	out := new([2]*rat.Rat)
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
	return *out
}

// Get one rat from each of 2 demand channels

func get2(in [2]*dch) [2]*rat.Rat {
	return getn([2]*dch{in[0], in[1]})
}

func copy(in *dch, out *dch) {
	for {
		<-out.req
		out.dat <- get(in)
	}
}

func repeat(dat *rat.Rat, out *dch) {
	for {
		put(dat, out)
	}
}
