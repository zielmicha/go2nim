# Compatibility with Nim types
import collections/interfaces, collections/goslice, collections/gcptrs, golib/builtin

proc `+`*(a: string, b: string): string =
  return a & b

proc `+=`*(a: var string, b: string) =
  a &= b

proc `[]`*(s: string, i: int64): char =
  return s[i.int]

proc toSlice*[T](s: seq[T]): GoSlice[T] =
  let slice = make(GoSlice[T], s.len)
  for i in 0..<s.len:
    slice[i] = s[i]
  return slice

converter toRune*(c: uint32): int32 =
  return c.int32

converter toRune*(c: int32): uint32 =
  return c.uint32

template trivialConverter(src, dst) =
  converter trivial*(c: src): dst =
    return c.dst

# Char conversions

trivialConverter char, uint8
trivialConverter char, uint16
trivialConverter char, uint32
trivialConverter char, uint64
trivialConverter char, uint
trivialConverter char, int8
trivialConverter char, int16
trivialConverter char, int32
trivialConverter char, int64
trivialConverter char, int

# Extending conversions

trivialConverter uint8, int
trivialConverter uint8, int32
trivialConverter uint8, int64

trivialConverter uint16, int
trivialConverter uint16, int32
trivialConverter uint16, int64

trivialConverter uint32, int64

# Float conversions

trivialConverter int64, float64
trivialConverter int, float64
trivialConverter int, float32

# ....

proc `<`*(a: int, b: byte): bool =
  return a < b.int

proc `<`*(a: int, b: char): bool =
  return a < b.int

# Nim bug? proc `shl`[T: SomeUnsignedInt](x, y: T): T doesn't match
proc `shl`*(a, b: uint): uint =
  return  `shl`[uint](a, b)

proc `shl`*(a: int, b: uint): int =
  return (`shl`(a.uint, b.uint)).int

proc explicitConvert*(a: GoSlice[byte], to: typedesc[string]): string =
  var s: string = newString(a.len)
  copyMem(addr s[0], addr a[0], a.len)
  return s

proc explicitConvert*(a: SomeInteger, to: typedesc[string]): string =
  panic("int -> string")

proc `-`*(x: uint64): uint64 =
  return 0.uint64 - x

proc `-`*(x: char, y: char): char =
  return (x.byte - y.byte).char

proc `+=`*(x: var char, y: char) =
  x = (x.byte + y.byte).char

proc `+=`*(x: var int32, y: int) =
  x = (x + y.int32).int32

proc `==`*(x: char, y: byte): bool =
  return x.byte == x

proc `==`*(x: char, y: Rune): bool =
  return x.Rune == y

proc `==`*(x: char, y: int): bool =
  return x.int == y

# TODO: remove these methods and fix make((INTLITERAL), ...)
trivialConverter int, uint32
trivialConverter int, int32
trivialConverter int32, uint16

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

converter narrowInt*(c: int): uint64 =
  return c.uint64

converter narrowInt*(c: int64): uint64 =
  return c.uint64

converter narrowInt*(c: int): uint =
  return c.uint

converter narrowInt*(c: uint): int =
  return c.int

# bug: 1.uint + 1 doesn't work

proc `+`*(a: int, b: uint): uint =
  return system.`+`(a.uint, b)

proc `+`*(a: int, b: uint64): uint64 =
  return a.uint64 + b

proc `+`*(a: int, b: char): int =
  return a + b.int

proc `go/`*(a: int, b: uint64): uint64 =
  return a.uint64 div b

proc `<`*(a: int, b: SomeUnsignedInt): bool =
  if a < 0:
    return true
  else:
    return a.uint64 < b.uint64

proc `<=`*(a: int, b: SomeUnsignedInt): bool =
  if a < 0:
    return true
  else:
    return a.uint64 <= b.uint64

proc `<=`*(a: SomeUnsignedInt, b: SomeUnsignedInt): bool =
  return system.`<=`(a, b)

proc `<=`*(a: SomeUnsignedInt, b: int): bool =
  if b < 0:
    return false
  else:
    return a.uint64 <= b.uint64

proc `<=`*(a: byte, b: byte): bool =
  return system.`<=`(a, b)

proc `or`*(a: int, b: uint8): uint8 =
  return a.uint8 or b

proc `+=`*(a: var uint8, b: int) =
  a = a + b.uint8

proc `+=`*(a: var int, b: int) =
  a = a + b

proc `+=`*(a: var uint, b: int) =
  a = a + b

proc `+=`*(a: var char, b: int) =
  a = (a.int + b).char

proc `*`*(a: int, b: SomeUnsignedInt): SomeUnsignedInt =
  return a.SomeUnsignedInt * b

proc `[]`*(a: string, i: uint64): char =
  return a[i.int]

proc `[]`*(a: gcptr, i: SomeInteger): auto =
  return (a[])[i]

converter toComplex*(a: int): complex128 =
  return (a.float64, 0.float64)

proc complex*(a: float64, b: float64): complex128 =
  return (a, b)

# ---

proc `<`*(a: int, b: int64): bool =
  return a.int64 < b

proc `*`*(a: int, b: uint64): uint64 =
  return a.uint64 * b.uint64

proc `$`*(a: uint): string =
  return $(a.uint64)

proc `$`*(a: uint8): string =
  return $(a.uint64)

proc `$`*(a: uint16): string =
  return $(a.uint64)

proc `$`*(a: uint32): string =
  return $(a.uint64)

# Nim issue #4543
# also, we {.overflowChecks: off.}

template goSupportMakeMod*(typ) =
  proc `mod`*(a: typ, b: typ): typ =
    if b == 0:
      raise newException(DivByZeroError, "unsigned division by zero")
    when a is SomeSignedInt:
      if b == -1 and a == a.low: # arbitrary overflow
        return 0

    return system.`mod`(a, b)
