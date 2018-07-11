// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dch // import "github.com/GoLangsam/powser/dch.big"

import (
	"math/big"
	// "github.com/GoLangsam/powser/big"
)

// ===========================================================================
// Beg of demand channel object

// Dch is a
// demand channel.
type Dch struct {
	ch  chan *big.Rat
	req chan struct{}
}

// New returns
// a (pointer to a) fresh
// unbuffered
// demand channel.
func New() *Dch {
	d := Dch{
		ch:  make(chan *big.Rat),
		req: make(chan struct{}),
	}
	return &d
}

// DchMakeBuff returns
// a (pointer to a) fresh
// buffered (with capacity=`cap`)
// mode channel.
func DchMakeBuff(cap int) *Dch {
	d := Dch{
		ch:  make(chan *big.Rat, cap),
		req: make(chan struct{}),
	}
	return &d
}

// Get is the comma-ok multi-valued form to receive and
// reports whether a received value was sent before the *big.Rat channel was closed.
//
// Get blocks until the requst is accepted and value `val` has been received from `from`.
func (from *Dch) Get() (val *big.Rat, open bool) {
	from.req <- struct{}{}
	val, open = <-from.ch
	return
}

// Drop is to be called by a consumer when finished requesting.
// The request channel is closed in order to broadcast this.
func (from *Dch) Drop() {
	close(from.req)
}

// From returns the handshaking channels
// (for use in `select` statements)
// to receive values:
//  `req` to send a request `req <- struct{}{}` and
//  `rcv` to reveive such requested value from.
func (from *Dch) From() (req chan<- struct{}, rcv <-chan *big.Rat) {
	return from.req, from.ch
}

// Next is the request method.
// It returns when a request was received
// and reports iff the request channel was open.
//
// Next blocks until a requsted is received.
//
// A sucessful Next is to be followed by one Send(v).
func (into *Dch) Next() bool {
	_, ok := <-into.req
	return ok
}

// Send is to be used after a successful Next()
func (into *Dch) Send(val *big.Rat) {
	into.ch <- val
}

// Put is the send-upon-request method
// - aka "myAnyChan <- myAny".
//
// Put is a convenience for
//  if Req() { Snd(v) }
//
// Put blocks until requsted to send value `val` into `into`.
func (into *Dch) Put(val *big.Rat) bool {
	_, ok := <-into.req
	if ok {
		into.ch <- val
	}
	return ok
}

// Into returns the handshaking channels
// (for use in `select` statements)
// to send values:
//  `req` to receive a request `<-req` and
//  `snd` to send such requested value into.
func (into *Dch) Into() (req <-chan struct{}, snd chan<- *big.Rat) {
	return into.req, into.ch
}

// Close is to be called by a producer when finished sending.
// The *big.Rat channel is closed in order to broadcast this.
//
// In order to avoid deadlock, pending requests are drained.
func (into *Dch) Close() {
	close(into.ch)
	for _ = range into.req {
	} // drain requests - there could be some
}

// Cap reports the capacity of the underlying *big.Rat channel.
func (c *Dch) Cap() int {
	return cap(c.ch)
}

// Len reports the length of the underlying *big.Rat channel.
func (c *Dch) Len() int {
	return len(c.ch)
}

// End of demand channel object
// ===========================================================================
