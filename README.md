# powser
Power series (with rational coefficients) by [lazy evaluated](https://en.wikipedia.org/wiki/Lazy_evaluation) demand channels. No goroutines leak!

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/GoLangsam/pipe)](https://goreportcard.com/report/github.com/GoLangsam/powser)
[![Build Status](https://travis-ci.org/GoLangsam/powser.svg?branch=master)](https://travis-ci.org/GoLangsam/powser)
[![GoDoc](https://godoc.org/github.com/GoLangsam/powser?status.svg)](https://godoc.org/github.com/GoLangsam/powser)

## Overview
[Power series](https://en.wikipedia.org/wiki/Formal_power_series) - with rational coefficients,

As _M. Douglas McIlroy_ says in his [Squinting at Power Series](https://swtch.com/~rsc/thread/squint.pdf):

> Data streams are an ideal vehicle for handling power series. Stream
> implementations can be read off directly from simple recursive equations
> that define operations such as multiplication, substitution, exponentiation,
> and reversion of series. The bookkeeping that bedevils these algorithms
> when they are expressed in traditional languages is completely hidden
> when they are expressed in stream terms. Communicating processes are
> the key to the simplicity of the algorithms.

`powser` builds on advanced concurrent functionality:
[lazy evaluated](https://en.wikipedia.org/wiki/Lazy_evaluation) demand channels
represent the potentially infinite stream of coefficients.

Coefficients can be `*big.Rat` from "math/big" or simple int64 rationals (included).

Powser has simple [conventions](#conventions) and provides:

- [Arithmetic](#arithmetic)
- [Constructors](#constructors)
- [Series](#series)
- [Consumers](#consumers)
- [Types](#types)
- [Helpers](#helpers)
- [Coefficients](#coefficients)

## Conventions
- Upper-case for power series.
- Lower-case for coefficients.
- Input variables:  From U,V,...
- Output variables: Into ...,Y,Z

## Remarks

Note: Use of coefficients from any other ring (such as square matrices) is easy to accomplish.
It just takes minimal dedicated changes - given the arithmetic methods are homo morph.

Note: `powser` is inspired by a test from the standard Go distribution (`test/chan/powser1.go`).

This implementation makes a clear separation between the power series themselves,
and the coefficients and their demand channel, which live in separate packages.

And all the spawned goroutines do not leak! They all terminate, both
- upon input being closed by the producer (lacking more data: a finite power series aka polynom) - and 
- upon output being dropped by the consumer (as there is no need for further coefficients).

Great care was taken to get this right (which is not trivial).

The overall result is too good, powerful and complete
as to serve as an example only in the related project [`pipe/s`](https://github.com/GoLangsam/pipe)
and thus is given here as an independent and stand-alone repository in the hope to be useful for someone.

---
## Details

### [Arithmetic](arithmetic.go)
provides methods to combine power series:

- basic helpers
  - `CMul` multiplies `U` by a constant `c` and returns `c*U`.
  - `MonMul` multiplies `U` by the monomial `x^n` and returns `x^n * U`.
  - `XMul` multiplies `U` by `x` (by the monomial `x^1`) and returns `x * U`.
  - `Shift` returns `c + x*U`

- Variadic algebraic operations:
  - `Plus`, `Less` and `Times`

- Substitution:
  - `MonSubst`: Monomial Substitution: `U(c*x^n)` Each Ui is multiplied by `c^i` and followed by n-1 zeros.
  - `Subst`: Substitute V for x in U: `U(V(x))` (also called "composition").

- The operators from differential calculus / analysis:
  - `Deriv`: Differentiates `U` and returns the derivative.
  - `Integ`: Integrate, with const of integration.

- further
  - `Recip`: the reciprocal of a power series: `1 / U`
  - `Exp` exponentiation of a power series (with constant term equal zero).


### [Constructors](constructors.go)
provides power series of given coefficient(s):

- Monomial: `c * x^n`
- Binomial: `(1+x)^c`, a finite polynom iff `c` is a positive and an alternating infinite power series otherwise.
- Polynom(a ...): converts coefficients, constant term first, to a (finite) power series, the polynom in the coefficients.

### [Series](series.go)
provides functions which return specific power series such as:

- `Factorials` starting from zero: 1, 1, 2, 6, 24, 120, 720, 5040, 40.320 ...
- `OneByFactorial` starting from zero: 1/1, 1/1, 1/2, 1/6, 1/120, 1/720 ...
- `Fibonaccis` starting from zero: 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 114 ...
- `OneByFibonacci` starting from zero: 1/1, 1/2, 1/3, 1/5, 1/8, 1/13, 1/21 ...
- `Sincos` returns the power series for sine and cosine (in radians).
- `Sin` returns the power series for sine (in radians).
- `Cos` returns the power series for cosine (in radians).

### [Consumers](consumers.go)
provides methods to evaluate or print a given power series such as:

- `EvalAt` evaluates a power series at `x=c` for up to `n` terms.
- `EvalN` evaluates a power series at `x=c` for up to `n` terms in floating point.
- `Printn` prints up to n terms of a power series.
- `Printer` returns a copy of `U`, and concurrently prints up to n terms of it. Useful to inspect formulas as it can be chained.
- `Print` one billion terms. Use at Your own risk ;-)

### [Types](types.go)
provides the bridge between the types used within the package (`Coefficient` and `PS`) and their imported realization.

As of now, there are `types.go.rat` and `types.go.math.big`.
Whichever is copied to `types.go` determines the chosen realization.

### [Helpers](helpers)
There are [wrappers](dch-wrap.go) to the underlying demand channel package:
- `Split` returns a pair of power series identical to the given one.
- `Append` all coefficients from `U` into `Into`.
- `GetNextFrom` `U` for `Into` and report success. Follow with `Into.Send( f(c) )`, iff ok.
- `GetWith` returns each first value received from the two given power series together with their respective ok boolean.
- `SendCfnFrom` applies a function `cfn(From)` to a coefficent received from `From`, sends the result into `Into` and report success.

And for the latter one there are closure functions on some coefficient math - for convenient use with `SendCfnFrom`

Note: Such closures are used where it helps to tighten the implementation of an algorithm,
and in other places calculations are intentionally done directly and explicit.

### [Coefficients](coefficients.go)
provides special (and non-exported) rational coefficients such a `aZero` or `aOne`.

---
