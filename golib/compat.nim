# Compatibility with Nim types
import collections/interfaces, collections/goslice, collections/gcptrs

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

converter toRune*(c: char): int32 =
  return c.int32

converter toRune*(c: uint32): int32 =
  return c.int32

converter toRune*(c: int32): uint32 =
  return c.uint32

converter toRune*(c: uint8): int32 =
  return c.int32

converter toByte*(c: char): byte =
  return c.byte

converter toInt*(c: uint16): int32 =
  return c.int32

converter toInt*(c: uint8): int =
  return c.int

converter toInt*(c: char): int =
  return c.int

proc `<`*(a: int, b: byte): bool =
  return a < b.int

proc `<`*(a: int, b: char): bool =
  return a < b.int

# TODO: remove these methods and fix make((INTLITERAL), ...)
converter narrowInt*(a: int): uint32 =
  return a.uint32

converter narrowInt*(a: int32): uint16 =
  return a.uint16

converter narrowInt*(a: uint32): uint8 =
  return a.uint8

converter narrowInt*(a: int32): uint8 =
  return a.uint8

converter narrowInt*(a: int): uint16 =
  return a.uint16

converter narrowInt*(c: int): uint8 =
  return c.uint8

converter narrowInt*(c: int32): uint16 =
  return c.uint8
