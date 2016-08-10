import collections/interfaces, collections/exttype, collections/gcptrs, collections/goslice, collections/macrotool
import macros, typetraits

exttypes:
  type
    Error = iface((
      error(): string
    ))

type
  Rune* = int32
  uintptr* = int

type
  Map*[K, V] = ref object
    # TODO
    a: int

  GoAutoArray*[T] = object

include builtincomplex
include builtinerrors

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
  when R is tuple[]:
    return T()
  elif R is tuple:
    makeObjectFromTuple(fields, t)
  else:
    makeObjectFromSingleItem(fields, t)

proc make*[T, R](fields: R, t: typedesc[ref T]): gcptr[T] =
  let ret = new(T)
  ret[] = make(fields, T)
  return ret

macro len(a: tuple): expr =
  let typ = a.getType
  if $typ[0].symbol != "tuple":
    error("invalid type")
  return newIntLitNode(typ.len - 1)

proc tupleToArray(a: tuple): auto =
  var s: seq[type(a[0])] = @[]
  for k, v in a.fieldPairs:
    s.add v
  return s

proc tupleToArray(a: array | seq): auto =
  return a

proc make*[T; R](fields: R, t: typedesc[GoAutoArray[T]]): auto =
  var ret: GoArray[T, fields.len]
  var fieldsArr = tupleToArray(fields)
  for i in 0..<fieldsArr.len:
    when fieldsArr[i] is tuple:
      let (k, v) = fieldsArr[i]
      ret[k] = v
    else:
      ret[i] = fieldsArr[i]
  return ret

proc make*[T; R; N: static[int32]](fields: R, t: typedesc[GoArray[T, N]]): GoArray[T, N] =
  when fields is tuple[]:
    return
  else:
    var fieldsArr = tupleToArray(fields)
    for i in 0..<fieldsArr.len:
      when fieldsArr[i] is tuple:
        let (k, v) = fieldsArr[i]
        result[k] = v
      else:
        result[i] = fieldsArr[i]

proc make*[K, V, R](fields: R, t: typedesc[Map[K, V]]): Map[K, V] =
  panic("make gomap")

proc make*[T, R](fields: R, t: typedesc[GoSlice[T]]): GoSlice[T] =
  var arr = make(fields, GoAutoArray[T])
  return arrToSlice(arr)

proc make*[T](t: typedesc[GoSlice[T]]): GoSlice[T] =
  panic("make goslice")

proc make*[K, V](t: typedesc[Map[K, V]]): Map[K, V] =
  panic("make gomap")

proc makeNilPtr*[T](t: typedesc[gcptr[T]]): gcptr[T] =
  return makeGcptr[T](nil, nil)

proc convert*[T, R](t: typedesc[T], v: R): T =
  when T is SomeInteger and R is SomeInteger:
    when nimvm: # TODO: nimvm doesn't support casts
      return T(v)
    else:
      return cast[T](v)
  elif compiles(T(v)):
    return T(v)
  elif compiles(explicitConvert(v, T)):
    return explicitConvert(v, T)
  else:
    when T is gcptr and R is NullType:
      return makeNilPtr(T)
    else:
      return v

proc maybeCastInterface*[T, R](v: T, to: typedesc[R]): tuple[val: R, ok: bool] =
  when T is not Interface:
    error("source type has to be an interface")

  if typeId(v) != typeId(R):
    result.ok = false
    return

  result.ok = true
  v.unwrapTo(addr result.val)

proc castInterface*[T, R](v: T, to: typedesc[R]): R =
  let (ret, ok) = maybeCastInterface(v, R)

  if not ok:
    panic("invalid cast (from " & typeId(v).name & " to " & typeId(R).name & ")")

  return ret

proc println*(args: varargs[string, `$`]) =
  echo(@args)

proc print*(args: varargs[string, `$`]) =
  echo(@args)

proc `go/`*(a: SomeInteger, b: SomeInteger): SomeInteger =
  if b == 0:
    panic("division by zero") # TODO: catch SIGFPE?
  when a is SomeSignedInt:
    if b == -1 and a == a.low: # arbitrary overflow
      return a
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

# proc append*[T](slice: GoSlice[T], elements: varargs[T]): GoSlice[T] =
#   return append(slice, elements)

proc append*[T](slice: GoSlice[T], a: T): GoSlice[T] =
  return append(slice, [a].make(GoSlice[T]))

proc append*[T](slice: GoSlice[T], a: T, b: T): GoSlice[T] =
  return append(slice, [a, b].make(GoSlice[T]))

proc append*[T](slice: GoSlice[T], a: T, b: T, c: T): GoSlice[T] =
  return append(slice, [a, b, c].make(GoSlice[T]))

type GoVarArgs*[T] = distinct GoSlice[T]

macro gofunc*(def: untyped): stmt =
  let pragmas = def[4]
  let params = def[3]

  pragmas.add(newIdentNode("discardable"))
  pragmas.add(newIdentNode("procvar"))

  if params.len > 1:
    # varargs handling
    let varargsDef = def.copyNimTree
    let lastArg = varargsDef[3][params.len - 1]
    if lastArg[1].kind == nnkBracketExpr and lastArg[1][0] == newIdentNode("GoVarArgs"):
      let varargType = lastArg[1][1]
      result = newNimNode(nnkStmtList).add(def)

      # wrapper using varargs[T]
      lastArg[1][0] = newIdentNode("varargs")
      var call = newCall(macrotool.stripPublic(def[0]))

      var i = 0
      for param in params:
        if param.kind == nnkIdentDefs:
          if i == len(params) - 1:
            call.add(newCall(newIdentNode("govarargs"), newCall(
              newIdentNode("toSlice"), newCall(newIdentNode("@"), param[0]))))
          else:
            call.add(param[0])
        i += 1

      if params[0] != newIdentNode("void"):
        call = newNimNode(nnkReturnStmt).add(call)

      result.add(varargsDef)
      varArgsDef[6] = call

      # wrapper by repeating args, disabled for now

      #[
      for j in 1..3:
        var call = newCall(macrotool.stripPublic(def[0]))
        var additionalNames: seq[NimNode] = @[]

        for param in params:
          if param.kind == nnkIdentDefs:
            call.add(param[0])

        call.del(call.len - 1)
        let callArgs = newNimNode(nnkBracket)
        call.add(callArgs)

        for i in 0..<j:
          let name = gensym(nskParam, ident="arg" & $i & "_")
          additionalNames.add(name)
          callArgs.add(newCall("convert", varargType, name))

        if params[0] != newIdentNode("void"):
          call = newNimNode(nnkReturnStmt).add(call)

        let manual = def.copyNimTree
        let manualParams = manual[3]
        manualParams.del(manualParams.len - 1)
        for name in additionalNames:
          manualParams.add(newNimNode(nnkIdentDefs).add(name, varargType, newNimNode(nnkEmpty)))
        manual[6] = call
        manual.repr.echo
        result.add(manual)
      #]

      return

  return def

proc govarargs*[T](a: GoSlice[T]): GoVarArgs[T] =
  return GoVarArgs[T](a)

converter toSlice*[T](a: GoVarArgs[T]): GoSlice[T] =
  return GoSlice[T](a)

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
  return newNimNode(nnkStmtList).add(defCopy, def)

macro runelit*(s: untyped): expr =
  if s.kind == nnkIntLit:
    return s
  elif s.kind == nnkCharLit:
    return newIntLitNode(s.intVal)
  else:
    error("bad runelit")

proc `==`*(a: proc, b: NullType): bool =
  return a == nil

when isMainModule:
  type
    Foo = object
      a: int
      b: string

  echo make((1, "2"), Foo)
