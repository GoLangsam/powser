// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package main

import (
	"fmt"
	"os"
	"time"
)

// ===========================================================================

func main() {

	if n > 0 {
		sample(n)
	} else {
		for i := 1; i <= max; i++ {
			sample(i)
		}
	}

	if x {
		fmt.Println("about to leave ...")
		<-time.After(time.Millisecond * 100)
		os.Exit(1) // to see leaking goroutines, if any
	}

}

// ===========================================================================
