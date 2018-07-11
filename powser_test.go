// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test concurrency primitives: power series.

package ps

import (
	"fmt"
	"os"
)

var (
	Ones PS
	Twos PS
)

func init() {
	Ones = AdInfinitum(aOne)
	Twos = AdInfinitum(aTwo)
}

// eq discriminates iff x is equal to y.
func eq(x, y Coefficient) bool {
	return x.Num() == y.Num() && y.Denom() == y.Denom()
}

func check(U PS, c Coefficient, count int, str string) {
	for i := 0; i < count; i++ {
		r, _ := U.Get()
		if !eq(r, c) {
			fmt.Print("got: ")
			fmt.Print(r.String())
			fmt.Print("should get ")
			fmt.Print(c.String())
			fmt.Print("\n")
			panic(str)
		}
	}
}

const N = 10

func checka(U PS, a []Coefficient, str string) {
	for i := 0; i < N; i++ {
		check(U, a[i], 1, str)
	}
}

func Example() {

	if len(os.Args) > 1 { // print
		fmt.Print("Ones: ")
		Ones.Printn(10)
		fmt.Print("Twos: ")
		Twos.Printn(10)
		fmt.Print("Add: ")
		Add(Ones, Twos).Printn(10)
		fmt.Print("Diff: ")
		Ones.Diff().Printn(10)
		fmt.Print("Integ: ")
		Ones.Integ(aZero).Printn(10)
		fmt.Print("CMul: ")
		Ones.Cmul(aMinusOne).Printn(10)
		fmt.Print("Sub: ")
		Ones.Minus(Twos).Printn(10)
		fmt.Print("Mul: ")
		Mul(Ones, Ones).Printn(10)
		fmt.Print("Exp: ")
		Ones.Exp().Printn(15)
		fmt.Print("MonSubst: ")
		Ones.MonSubst(aMinusOne, 2).Printn(10)
		fmt.Print("ATan: ")
		Ones.MonSubst(aMinusOne, 2).Integ(aZero).Printn(10)
	} else { // test
		check(Ones, aOne, 5, "Ones")
		check(Add(Ones, Ones), aTwo, 0, "Add Ones Ones")                 // 1 1 1 1 1
		check(Add(Ones, Twos), NewCoefficient(3, 1), 0, "Add Ones Twos") // 3 3 3 3 3
		a := make([]Coefficient, N)
		d := Ones.Diff()
		for i := 0; i < N; i++ {
			a[i] = NewCoefficient(int64(i+1), 1)
		}
		checka(d, a, "Diff") // 1 2 3 4 5
		in := Ones.Integ(aZero)
		a[0] = aZero // integration constant
		for i := 1; i < N; i++ {
			a[i] = NewCoefficient(1, int64(i))
		}
		checka(in, a, "Integ")                                         // 0 1 1/2 1/3 1/4 1/5
		check(Twos.Cmul(aMinusOne), NewCoefficient(-2, 1), 10, "CMul") // -1 -1 -1 -1 -1
		check(Ones.Minus(Twos), aMinusOne, 0, "Sub Ones Twos")         // -1 -1 -1 -1 -1
		m := Mul(Ones, Ones)
		for i := 0; i < N; i++ {
			a[i] = NewCoefficient(int64(i+1), 1)
		}
		checka(m, a, "Mul") // 1 2 3 4 5
		e := Ones.Exp()
		a[0] = aOne
		a[1] = aOne
		a[2] = NewCoefficient(3, 2)
		a[3] = NewCoefficient(13, 6)
		a[4] = NewCoefficient(73, 24)
		a[5] = NewCoefficient(167, 40)
		a[6] = NewCoefficient(4051, 720)
		a[7] = NewCoefficient(37633, 5040)
		a[8] = NewCoefficient(43817, 4480)
		a[9] = NewCoefficient(4596553, 362880)
		checka(e, a, "Exp") // 1 1 3/2 13/6 73/24
		at := Ones.MonSubst(aMinusOne, 2).Integ(aZero)
		for c, i := 1, 0; i < N; i++ {
			if i%2 == 0 {
				a[i] = aZero
			} else {
				a[i] = NewCoefficient(int64(c), int64(i))
				c *= -1
			}
		}
		checka(at, a, "ATan") // 0 -1 0 -1/3 0 -1/5
		/*
			t := Revert(Ones.MonSubst(aMinusOne, 2).Integ(aZero))
			a[0] = aZero
			a[1] = aOne
			a[2] = aZero
			a[3] = NewCoefficient(1,3)
			a[4] = aZero
			a[5] = NewCoefficient(2,15)
			a[6] = aZero
			a[7] = NewCoefficient(17,315)
			a[8] = aZero
			a[9] = NewCoefficient(62,2835)
			checka(t, a, "Tan")  // 0 1 0 1/3 0 2/15
		*/
	}
}
