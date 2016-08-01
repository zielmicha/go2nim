import complex

type
  complex128* = Complex
  complex64* = complex128 # TODO

export complex

proc real*(a: complex128): float64 =
  return a.re

proc imag*(a: complex128): float64 =
  return a.im
