# Translational module
import gomathwrapper
export gomathwrapper

proc inf*(a: int): float64 =
  if a >= 0:
    return 1e1000
  else:
    return -1e1000

proc naN*(): float64 =
  return inf(1) + inf(-1)

proc float32frombits*(b: uint32): float32 =
  return cast[float32](b)

proc float32bits*(b: float32): uint32 =
  return cast[uint32](b)

proc float64bits*(b: float64): uint64 =
  return cast[uint64](b)

proc float64frombits*(b: uint64): float64 =
  return cast[float64](b)
