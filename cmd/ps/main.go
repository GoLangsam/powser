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

	doTests()

	if x {
		fmt.Println("about to leave ... are goroutines leaking?")
		<-time.After(time.Millisecond * 100)
		os.Exit(1) // to see leaking goroutines, if any
	}

}

// ===========================================================================
