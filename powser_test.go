// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package ps_test

import (
	"testing"

	. "github.com/GoLangsam/powser"
)

func check(t *testing.T, U PS, c Coefficient, count int, str string) {
	chk(t, U, c, count, str)
	U.Drop()
}

func chk(t *testing.T, U PS, c Coefficient, count int, str string) {
	for i := 0; i < count; i++ {

		if r, ok := U.Receive(); ok {
			if !Equal(r, c) {
				t.Error("got: ", r.String(), "\t should get ", c.String(), "\t ", str)
			}
		}
	}
}

const N = 10

func checka(t *testing.T, U PS, a []Coefficient, str string) {
	for i := 0; i < len(a); i++ {
		chk(t, U, a[i], 1, str)
	}
	U.Drop()
}

func TestPS(t *testing.T) {

	check(t, Ones(), aOne(), 5, "Ones()")
	check(t, Ones().Plus(Ones()), NewCoefficient(2, 1), 5, "Add Ones() Ones()") // 1 1 1 1 1
	check(t, Ones().Plus(Twos()), NewCoefficient(3, 1), 5, "Add Ones() Twos()") // 3 3 3 3 3

	a := make([]Coefficient, N)
	d := Ones().Deriv()
	for i := 0; i < N; i++ {
		a[i] = ratIby1(i + 1)
	}
	checka(t, d, a, "Deriv") // 1 2 3 4 5

	in := Ones().Integ(aZero())
	a[0] = aZero() // integration constant
	for i := 1; i < N; i++ {
		a[i] = rat1byI(i)
	}
	checka(t, in, a, "Integ") // 0 1 1/2 1/3 1/4 1/5

	check(t, Twos().CMul(aMinusOne()), NewCoefficient(-2, 1), 10, "CMul") // -1 -1 -1 -1 -1
	check(t, Ones().Less(Twos()), aMinusOne(), 5, "Sub Ones() Twos()")    // -1 -1 -1 -1 -1

	m := Ones().Times(Ones())
	for i := 0; i < N; i++ {
		a[i] = ratIby1(i + 1)
	}
	checka(t, m, a, "Mul") // 1 2 3 4 5

	e := Ones().Exp()
	a[0] = aOne()
	a[1] = aOne()
	a[2] = NewCoefficient(3, 2)
	a[3] = NewCoefficient(13, 6)
	a[4] = NewCoefficient(73, 24)
	a[5] = NewCoefficient(167, 40)
	a[6] = NewCoefficient(4051, 720)
	a[7] = NewCoefficient(37633, 5040)
	a[8] = NewCoefficient(43817, 4480)
	a[9] = NewCoefficient(4596553, 362880)
	checka(t, e, a, "Exp") // 1 1 3/2 13/6 73/24

	at := Ones().MonSubst(aMinusOne(), 2).Integ(aZero())
	for c, i := 1, 0; i < N; i++ {
		if i%2 == 0 {
			a[i] = aZero()
		} else {
			a[i] = NewCoefficient(int64(c), int64(i))
			c *= -1
		}
	}
	checka(t, at, a, "ATan") // 0 1 0 -1/3 0 1/5 0 -1/7 0 1/9 0 -1/11

	/*
		t := Revert(Ones().MonSubst(aMinusOne(), 2).Integ(aZero()))
		a[0] = aZero()
		a[1] = aOne()
		a[2] = aZero()
		a[3] = NewCoefficient(1,3)
		a[4] = aZero()
		a[5] = NewCoefficient(2,15)
		a[6] = aZero()
		a[7] = NewCoefficient(17,315)
		a[8] = aZero()
		a[9] = NewCoefficient(62,2835)
		checka(t, t, a, "Tan")  // 0 1 0 1/3 0 2/15
	*/
}
