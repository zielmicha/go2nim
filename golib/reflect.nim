include gosupport

type
  Value* = object
    val: EmptyInterface

  Type* = ref object

type Kind* {.pure.} = enum
  invalid = 0
  bool = 1
  int = 2
  int8 = 3
  int16 = 4
  int32 = 5
  int64 = 6
  uint = 7
  uint8 = 8
  uint16 = 9
  uint32 = 10
  uint64 = 11
  uintptr = 12
  float32 = 13
  float64 = 14
  complex64 = 15
  complex128 = 16
  array = 17
  chan = 18
  `func` = 19
  `interface` = 20
  map = 21
  `ptr` = 22
  slice = 23
  string = 24
  struct = 25
  unsafePointer = 26

proc valueOf*(v: EmptyInterface): Value =
  return

proc field*(v: Value, index: int): Value =
  panic("not implemented")

proc kind*(v: Value): Kind =
  panic("not implemented")

proc isNil*(v: Value): bool =
  panic("not implemented")

proc elem*(v: Value): Value =
  ## Unpack interface or pointer.
  panic("not implemented")

proc isValid*(v: Value): bool =
  panic("not implemented")

proc getType*(v: Value): Type =
  return

proc typeOf*(v: EmptyInterface): Type =
  if v == null:
    return nil
  return valueOf(v).getType

proc canAddr*(v: Value): bool =
  panic("not implemented")

proc canInterface*(v: Value): bool =
  panic("not implemented")

proc `interface`*(v: Value): EmptyInterface =
  panic("not implemented")

proc toPointer*(v: Value): int =
  panic("not implemented")

proc toBool*(v: Value): bool =
  panic("not implemented")

proc toInt*(v: Value): int =
  panic("not implemented")

proc toUint*(v: Value): uint =
  panic("not implemented")

proc toFloat*(v: Value): float =
  panic("not implemented")

proc toString*(t: Value): string =
  panic("not implemented")

proc complex*(v: Value): complex128 =
  panic("not implemented")

proc bytes*(v: Value): GoSlice[byte] =
  panic("not implemented")

proc setBool*(v: Value, val: bool) =
  panic("not implemented")

proc setInt*(v: Value, val: int64) =
  panic("not implemented")

proc setUint*(v: Value, val: uint64) =
  panic("not implemented")

proc setString*(v: Value, val: string) =
  panic("not implemented")

proc setFloat*(v: Value, val: float64) =
  panic("not implemented")

proc setComplex*(v: Value, val: complex128) =
  panic("not implemented")

proc set*(v: Value, val: Value) =
  panic("not implemented")

proc makeSlice*(typ: Type, len: int, cap: int): Value =
  panic("not implemented")

proc mapKeys*(v: Value): GoSlice[Value] =
  panic("not implemented")

proc mapIndex*(v: Value, key: Value): Value =
  panic("not implemented")

proc len*(v: Value): int =
  panic("not implemented")

proc slice*(v: Value, a: int, b: int): Value =
  panic("not implemented")

proc index*(v: Value, a: int): Value =
  panic("not implemented")

proc numField*(v: Value): int =
  panic("not implemented")

proc field*(v: Type, index: int): Type =
  # TODO: this method does not appear in docs
  panic("not implemented")

proc bits*(v: Type): int =
  panic("not implemented")

proc elem*(v: Type): Type =
  panic("not implemented")

proc kind*(v: Type): Kind =
  panic("not implemented")

proc name*(v: Type): string =
  panic("not implemented")

proc toString*(t: Type): string =
  panic("not implemented")
