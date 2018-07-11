// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package ps

import (
	"os"

	rat "github.com/GoLangsam/powser/big"
)

func Init() {
	Ones = Rep(rat.One)
	Twos = Rep(rat.NewRat(2, 1))
}

func check(U PS, c *rat.Rat, count int, str string) {
	for i := 0; i < count; i++ {
		r := U.Get()
		if !r.Eq(c) {
			print("got: ")
			r.Pr()
			print("should get ")
			c.Pr()
			print("\n")
			panic(str)
		}
	}
}

const N = 10

func checka(U PS, a []*rat.Rat, str string) {
	for i := 0; i < N; i++ {
		check(U, a[i], 1, str)
	}
}

func Example() {
	Init()
	if len(os.Args) > 1 { // print
		print("Ones: ")
		Printn(Ones, 10)
		print("Twos: ")
		Printn(Twos, 10)
		print("Add: ")
		Printn(Add(Ones, Twos), 10)
		print("Diff: ")
		Printn(Diff(Ones), 10)
		print("Integ: ")
		Printn(Integ(rat.Zero, Ones), 10)
		print("CMul: ")
		Printn(Cmul(rat.MinusOne, Ones), 10)
		print("Sub: ")
		Printn(Sub(Ones, Twos), 10)
		print("Mul: ")
		Printn(Mul(Ones, Ones), 10)
		print("Exp: ")
		Printn(Exp(Ones), 15)
		print("MonSubst: ")
		Printn(MonSubst(Ones, rat.MinusOne, 2), 10)
		print("ATan: ")
		Printn(Integ(rat.Zero, MonSubst(Ones, rat.MinusOne, 2)), 10)
	} else { // test
		check(Ones, rat.One, 5, "Ones")
		check(Add(Ones, Ones), rat.Two, 0, "Add Ones Ones")          // 1 1 1 1 1
		check(Add(Ones, Twos), rat.NewRat(3, 1), 0, "Add Ones Twos") // 3 3 3 3 3
		a := make([]*rat.Rat, N)
		d := Diff(Ones)
		for i := 0; i < N; i++ {
			a[i] = rat.NewRat(int64(i+1), 1)
		}
		checka(d, a, "Diff") // 1 2 3 4 5
		in := Integ(rat.Zero, Ones)
		a[0] = rat.Zero // integration constant
		for i := 1; i < N; i++ {
			a[i] = rat.NewRat(1, int64(i))
		}
		checka(in, a, "Integ")                                         // 0 1 1/2 1/3 1/4 1/5
		check(Cmul(rat.MinusOne, Twos), rat.NewRat(-2, 1), 10, "CMul") // -1 -1 -1 -1 -1
		check(Sub(Ones, Twos), rat.MinusOne, 0, "Sub Ones Twos")       // -1 -1 -1 -1 -1
		m := Mul(Ones, Ones)
		for i := 0; i < N; i++ {
			a[i] = rat.NewRat(int64(i+1), 1)
		}
		checka(m, a, "Mul") // 1 2 3 4 5
		e := Exp(Ones)
		a[0] = rat.One
		a[1] = rat.One
		a[2] = rat.NewRat(3, 2)
		a[3] = rat.NewRat(13, 6)
		a[4] = rat.NewRat(73, 24)
		a[5] = rat.NewRat(167, 40)
		a[6] = rat.NewRat(4051, 720)
		a[7] = rat.NewRat(37633, 5040)
		a[8] = rat.NewRat(43817, 4480)
		a[9] = rat.NewRat(4596553, 362880)
		checka(e, a, "Exp") // 1 1 3/2 13/6 73/24
		at := Integ(rat.Zero, MonSubst(Ones, rat.MinusOne, 2))
		for c, i := 1, 0; i < N; i++ {
			if i%2 == 0 {
				a[i] = rat.Zero
			} else {
				a[i] = rat.NewRat(int64(c), int64(i))
				c *= -1
			}
		}
		checka(at, a, "ATan") // 0 -1 0 -1/3 0 -1/5
		/*
			t := Revert(Integ(rat.Zero, MonSubst(Ones, rat.MinusOne, 2)))
			a[0] = rat.Zero
			a[1] = rat.One
			a[2] = rat.Zero
			a[3] = rat.NewRat(1,3)
			a[4] = rat.Zero
			a[5] = rat.NewRat(2,15)
			a[6] = rat.Zero
			a[7] = rat.NewRat(17,315)
			a[8] = rat.Zero
			a[9] = rat.NewRat(62,2835)
			checka(t, a, "Tan")  // 0 1 0 1/3 0 2/15
		*/
	}
}
