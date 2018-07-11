// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package dch

import (
	"github.com/GoLangsam/powser/big"
)

// Dch represents a demand channel
type Dch struct {
	req chan struct{}
	dat chan *big.Rat
}

// New reurns a (pointer to a) fresh demand channel
func New() *Dch {
	d := new(Dch)
	d.req = make(chan struct{})
	d.dat = make(chan *big.Rat)
	return d
}

// Into returns the handshaking channels for send
// `req` to receive a request and
// `dat` to send data to.
// Intended for use in `select` statements.
func (into *Dch) Into() (req <-chan struct{}, dat chan<- *big.Rat) {
	return into.req, into.dat
}

// From returns the handshaking channels for receive
// `req` to send a request and
// `dat` to reveive data from.
// Intended for use in `select` statements.
func (from *Dch) From() (req chan<- struct{}, dat <-chan *big.Rat) {
	return from.req, from.dat
}

// Put blocks until requsted to send dat to out
func (into *Dch) Put(dat *big.Rat) {
	<-into.req
	into.dat <- dat
}

// Get blocks until the requsted data can be returned
func (from *Dch) Get() (dat *big.Rat) {
	from.req <- struct{}{}
	return <-from.dat
}

// Copy data from `from` into `into`
func (into *Dch) Copy(from *Dch) {
	for {
		<-into.req
		into.dat <- from.Get()
	}
}

// Repeat keeps sending `dat` into `into`
func (into *Dch) Repeat(dat *big.Rat) {
	for {
		into.Put(dat)
	}
}
