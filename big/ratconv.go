// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big

import (
	"fmt"
)

// String returns a string representation of x in the form "( a / b )" (or "a" iff b == 1).
func (x *Rat) String() string {
	if x.den == 1 {
		return fmt.Sprint(x.num)
	} else {
		return fmt.Sprint("( ", x.num, " / ", x.den, " )")
	}
}

func (u *Rat) Pr() {
	if u.den == 1 {
		print(u.num)
	} else {
		print(u.num, "/", u.den)
	}
	print(" ")
}
