// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dch // import "github.com/GoLangsam/powser/big.dch"

import (
	"math/big"
	// "github.com/GoLangsam/powser/rat"
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
// buffered
// demand channel
// (with capacity=`cap`).
func DchMakeBuff(cap int) *Dch {
	d := Dch{
		ch:  make(chan *big.Rat, cap),
		req: make(chan struct{}),
	}
	return &d
}

// ---------------------------------------------------------------------------

// Get is the comma-ok multi-valued form to receive from the channel and
// reports whether a received value was sent before the channel was closed.
//
// Get blocks until the request is accepted and value `val` has been received from `from`.
func (from *Dch) Get() (val *big.Rat, open bool) {
	from.req <- struct{}{}
	val, open = <-from.ch
	return
}

// Drop is to be called by a consumer when finished requesting.
// The request channel is closed in order to broadcast this.
//
// In order to avoid deadlock, pending sends are drained.
func (from *Dch) Drop() {
	close(from.req)
	go func(from *Dch) {
		for _ = range from.ch {
		} // drain values - there could be some
	}(from)
}

// From returns the handshaking channels
// (for use in `select` statements)
// to receive values:
//  `req` to send a request `req <- struct{}{}` and
//  `rcv` to reveive such requested value from.
func (from *Dch) From() (req chan<- struct{}, rcv <-chan *big.Rat) {
	return from.req, from.ch
}

// ---------------------------------------------------------------------------

// GetNextFrom `from` for `into` and report success.
// Follow it with `into.Send( f(val) )`, if ok.
func (into *Dch) GetNextFrom(from *Dch) (val *big.Rat, ok bool) {
	if ok = into.Next(); ok {
		val, ok = from.Get()
	}
	if !ok {
		from.Drop()
		into.Close()
	}
	return
}

// Next is the request method.
// It returns when a request was received
// and reports whether the request channel was open.
//
// Next blocks until a requested is received.
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
//  if Next() { Send(v) }
//
// Put blocks until requested to send value `val` into `into`.
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
// The value channel is closed in order to broadcast this.
//
// In order to avoid deadlock, pending requests are drained.
func (into *Dch) Close() {
	close(into.ch)
	go func(into *Dch) {
		for _ = range into.req {
		} // drain requests - there could be some
	}(into)
}

// ---------------------------------------------------------------------------

// MyDch returns itself.
func (c *Dch) MyDch() *Dch {
	return c
}

// Cap reports the capacity of the underlying value channel.
func (c *Dch) Cap() int {
	return cap(c.ch)
}

// Len reports the length of the underlying value channel.
func (c *Dch) Len() int {
	return len(c.ch)
}

// End of demand channel object
// ===========================================================================
