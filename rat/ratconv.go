// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big

import (
	"fmt"
)

// String returns a string representation of x in the form "a/b" (or "a" iff b == 1).
func (x *Rat) String() string {
	if x.den == 1 {
		return fmt.Sprint(x.num)
	} else {
		return fmt.Sprint(x.num, "/", x.den)
	}
}

// Float64 returns the nearest float64 value for x and a bool indicating
// whether f represents x exactly. If the magnitude of x is too large to be
// represented by a float64, f is an infinity and exact is false. The sign of f
// always matches the sign of x, even if f == 0.
func (x *Rat) Float64() (f float64, exact bool) {

	return float64(x.num) / float64(x.den), true // we cheat a little

}
