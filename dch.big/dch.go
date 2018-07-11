// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package dch // import "github.com/GoLangsam/powser/dch.big"

import (
	"math/big"
	// "github.com/GoLangsam/ps/big"
)

// Dch represents a demand channel.
type Dch struct {
	req chan struct{}
	dat chan *big.Rat
}

// New reurns a (pointer to a) fresh demand channel.
func New() *Dch {
	d := new(Dch)
	d.req = make(chan struct{})
	d.dat = make(chan *big.Rat)
	return d
}

// Into returns the handshaking channels
// (for use in `select` statements)
// to send values:
//  `req` to receive a request `<-req` and
//  `snd` to send such requested value into.
func (into *Dch) Into() (req <-chan struct{}, snd chan<- *big.Rat) {
	return into.req, into.dat
}

// From returns the handshaking channels
// (for use in `select` statements)
// to receive values:
//  `req` to send a request `req <- struct{}{}` and
//  `rcv` to reveive such requested value from.
func (from *Dch) From() (req chan<- struct{}, rcv <-chan *big.Rat) {
	return from.req, from.dat
}

// Put blocks until requsted to send value `val` into `into`.
func (into *Dch) Put(val *big.Rat) {
	<-into.req
	into.dat <- val
}

// Get blocks until the requst is accepted and value `val` has been received from `from`.
func (from *Dch) Get() (val *big.Rat, ok bool) {
	from.req <- struct{}{}
	val, ok = <-from.dat
	return
}

// Close closes the underlying channel
func (into *Dch) Close() {
	close(into.dat)
}
