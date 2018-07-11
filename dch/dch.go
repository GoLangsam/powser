// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package dch

import (
	"github.com/GoLangsam/powser/rat"
)

type Dch struct {
	req chan int
	dat chan *rat.Rat
	nam int
}

func (in *Dch) Req() <-chan int {
	return in.req
}

func (in *Dch) Dat() chan<- *rat.Rat {
	return in.dat
}

var chnames string
var chnameserial int
var seqno int

func init() {
	chnameserial = -1
	seqno = 0
	chnames = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
}

func New() *Dch {
	c := chnameserial % len(chnames)
	chnameserial++
	d := new(Dch)
	d.req = make(chan int)
	d.dat = make(chan *rat.Rat)
	d.nam = c
	return d
}

func (out *Dch) Put(dat *rat.Rat) {
	<-out.req
	out.dat <- dat
}

func (in *Dch) Get() *rat.Rat {
	seqno++
	in.req <- seqno
	return <-in.dat
}

func (out *Dch) Copy(in *Dch) {
	for {
		<-out.req
		out.dat <- in.Get()
	}
}

func (out *Dch) Repeat(dat *rat.Rat) {
	for {
		out.Put(dat)
	}
}
