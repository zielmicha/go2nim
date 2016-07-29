import collections/interfaces, collections/exttype, collections/gcptrs, collections/goslice, collections/macrotool
import macros

exttypes:
  type
    Error = iface((
      error(): string
    ))

type Rune* = int32

type uintptr* = int

type
  Map*[K, V] = ref object
    # TODO
    a: int

  GoAutoArray*[T] = object

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
  # objType.treeRepr.echo
  for field in objType[2]:
    ret.add(newNimNode(nnkExprColonExpr).add(field, newNimNode(nnkBracketExpr).add(fields, newIntLitNode(i))))
    i += 1

  # ret.repr.echo
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

proc make*[T; R](fields: R, t: typedesc[GoAutoArray[T]]): auto =
  var ret: GoArray[T, fields.len]
  panic("make goarray")
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
  when compiles(T(v)):
    return T(v)
  else:
    # necessary due to Nim bug?
    when T is gcptr and R is NullType:
      return makeNilPtr(T)
    else:
      return v

proc maybeCastInterface*[T, R](v: T, to: typedesc[R]): tuple[val: R, ok: bool] =
  return (convert(R, null), false) # TODO

proc castInterface*[T, R](v: T, to: typedesc[R]): R =
  let (ret, ok) = maybeCastInterface(v, R)
  if not ok:
    panic("invalid cast")

  return ret

proc `go/`*(a: SomeInteger, b: SomeInteger): SomeInteger =
  return a div b

proc `go/`*(a: float, b: float): float =
  return a / b

proc append*[T](slice: GoSlice[T], elements: GoSlice[T]): GoSlice[T] =
  var slice = slice
  var n = len(slice)
  var total = (len(slice) + len(elements))
  if (total > cap(slice)):
    var newSize = (((total * 3) div 2) + 1)
    var newSlice = make(GoSlice[T], total, newSize)
    copy(newSlice, slice)
    slice = newSlice
  slice = slice(slice, high=total)
  copy(slice(slice, low=n), elements)
  return slice

proc append*[T](slice: GoSlice[T], elements: varargs[T]): GoSlice[T] =
  return append(slice, elements)

proc append*[T](slice: GoSlice[T], a: T, b: T): GoSlice[T] =
  return append(slice, @[a, b])

proc append*[T](slice: GoSlice[T], a: T, b: T, c: T): GoSlice[T] =
  return append(slice, @[a, b, c])

macro gomethod*(def: untyped): stmt =
  # support for gcptr receivers
  # FIXME: THIS IS UNSAFE IF THE METHOD SAVES RECEIVER!
  let defCopy = def.copyNimTree
  let params = def[3]
  let firstParamType = params[1][1]

  var args: seq[NimNode] = @[]
  for param in params:
    if param.kind == nnkIdentDefs:
      args.add(param[0])

  var call = newCall(macrotool.stripPublic(def[0]), args)

  if params[0] != newIdentNode("void"):
    call = newNimNode(nnkReturnStmt).add(call)

  if firstParamType.kind == nnkBracketExpr and firstParamType[0].kind == nnkIdent and $firstParamType[0].ident == "gcptr":
    def[4] = newNimNode(nnkEmpty) # reset pragma

    params[1][1] = newNimNode(nnkVarTy).add(firstParamType[1])

    let selfName = params[1][0]
    def[6] = quote do:
      let `selfName` = makeGcPtr(addr `selfName`, nil)

  else:
    params[1][1] = newNimNode(nnkBracketExpr).add(newIdentNode("gcptr"), firstParamType)

    let selfName = params[1][0]
    def[6] = quote do:
      let `selfName` = `selfName`[]

  def[6].add(call)
  def.repr.echo
  return newNimNode(nnkStmtList).add(defCopy, def)

when isMainModule:
  type
    Foo = object
      a: int
      b: string

  echo make((1, "2"), Foo)
