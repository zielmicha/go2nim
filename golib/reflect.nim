include gosupport

type
  Value* = object
    val: EmptyInterface

  Type* = object

type Kind {.pure.} = enum
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

proc typeOf*(v: EmptyInterface): Type =
  panic("not implemented")

proc field*(v: Value, index: int): Value =
  panic("not implemented")

proc kind*(v: Value): Kind =
  panic("not implemented")
