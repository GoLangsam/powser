// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dch

// ===========================================================================

// Append synchronously from `from` into `into`.
func (into *Dch) Append(from *Dch) {
	defer func() {
		from.Drop()
		into.Close()
	}()
	into.AppendOnly(from)
}

// AppendOnly synchronously from `from` into `into`
// without cleanup of handshaking resources.
func (into *Dch) AppendOnly(from *Dch) {
	for into.Next() {
		if c, ok := from.Get(); ok {
			into.Send(c)
		} else {
			return
		}
	}
}

// ===========================================================================
