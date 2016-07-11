# Compatibility with Nim types
import collections/interfaces
import unicode

proc `+`*(a: string, b: string): string =
  return a & b

proc `+=`*(a: var string, b: string) =
  a &= b
