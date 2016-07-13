import collections/interfaces, collections/exttype, collections/gcptrs, collections/goslice
import macros

exttypes:
  type
    Error = iface((
      error(): string
    ))

type Rune* = int32

type
  Map*[K, V] = ref object
    # TODO
    a: int

proc panic*(err: string) =
  # TODO
  raise newException(Exception, "panic")

macro makeObjectFromTuple(fields: typed, t: typed): expr =
  let typeName = t.getType[1]
  var objType = typeName.getType
  if objType[0].kind == nnkSym:
    if $(objType[0].symbol) != "ref":
      error("unknown type")
    objType = objType[1].getType
  let ret = newNimNode(nnkObjConstr).add(typeName)
  var i = 0
  objType.treeRepr.echo
  for field in objType[2]:
    ret.add(newNimNode(nnkExprColonExpr).add(field, newNimNode(nnkBracketExpr).add(fields, newIntLitNode(i))))
    i += 1

  ret.repr.echo
  return ret

macro makeObjectFromSingleItem(item: typed, t: typed): expr =
  var typeName = t.getType[1]
  var objType = typeName.getType
  if objType[0].kind == nnkSym:
    if $(objType[0].symbol) != "ref":
      error("unknown type")
    objType = objType[1].getType
  let ret = newNimNode(nnkObjConstr).add(typeName)
  if objType[2].len != 1:
    error("object has more than one field")

  let field = objType[2][0]

  ret.add(newNimNode(nnkExprColonExpr).add(field, item))
  return ret

proc make*[T: object, R](fields: R, t: typedesc[T]): T =
  when R is tuple:
    makeObjectFromTuple(fields, t)
  else:
    makeObjectFromSingleItem(fields, t)

proc make*[T, R](fields: R, t: typedesc[ref T]): gcptr[T] =
  let ret = new(T)
  ret[] = make(fields, T)
  return ret

proc make*[T; R; N: static[int32]](fields: R, t: typedesc[GoArray[T, N]]): GoArray[T, N] =
  panic("make goarray")

proc make*[K, V, R](fields: R, t: typedesc[Map[K, V]]): Map[K, V] =
  panic("make gomap")

proc make*[T, R](fields: R, t: typedesc[GoSlice[T]]): GoSlice[T] =
  panic("make goslice")

proc make*[T](t: typedesc[GoSlice[T]]): GoSlice[T] =
  panic("make goslice")

proc make*[K, V](t: typedesc[Map[K, V]]): Map[K, V] =
  panic("make gomap")

proc makeNilPtr*[T](t: typedesc[gcptr[T]]): gcptr[T] =
  return makeGcptr[T](nil, nil)

proc convert*[T, R](t: typedesc[T], v: R): T =
  # necessary due to Nim bug?
  when T is gcptr and R is NullType:
    return makeNilPtr(T)
  else:
    return v

proc castInterface*[T, R](v: T, to: typedesc[R]): tuple[val: R, ok: bool] =
  return (convert(R, null), false) # TODO

proc `go/`*(a: SomeInteger, b: SomeInteger): SomeInteger =
  return a div b

proc `go/`*(a: float, b: float): float =
  return a / b

when isMainModule:
  type
    Foo = object
      a: int
      b: string

  echo make((1, "2"), Foo)
