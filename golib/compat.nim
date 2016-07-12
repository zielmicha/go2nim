# Compatibility with Nim types
import collections/interfaces, collections/goslice, collections/gcptrs
import unicode

proc `+`*(a: string, b: string): string =
  return a & b

proc `+=`*(a: var string, b: string) =
  a &= b

converter toSlice*[T](s: varargs[T]): GoSlice[T] =
  # FIXME: this causes problems with compile-time algorithm.sort
  let slice = make(GoSlice[T], s.len)
  for i in 0..s.len:
    slice[i] = s[i]
  return slice
