# Compatibility with Nim types

proc `+`*(a: string, b: string): string =
  return a & b

proc `+=`*(a: var string, b: string) =
  a &= b

