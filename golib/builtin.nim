import collections/interfaces, collections/exttype, collections/gcptrs
import macros

exttypes:
  type
    Error = iface((
      error(): string
    ))

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

proc panic*(err: string) =
  # TODO
  raise newException(Exception, "panic")

when isMainModule:
  type
    Foo = object
      a: int
      b: string

  echo make((1, "2"), Foo)
