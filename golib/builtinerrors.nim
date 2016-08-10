import sequtils

var exitDefer {.threadvar.}: bool

type
  PanicException* = object of Exception
    wrappedValue*: EmptyInterface

macro goPartialCall*(s: untyped): expr =
  result = newNimNode(nnkStmtList)
  var argNames: seq[NimNode] = @[]

  for arg in toSeq(s)[1..^1]:
    let argName = genSym(ident="arg")
    result.add(quote do:
      let `argName` = `arg`)

    argNames.add(argName)

  let callExpr = newNimNode(nnkCall)
  callExpr.add(s[0])
  for name in argNames:
    callExpr.add(name)

  result.add(newNimNode(nnkLambda).add(
    newNimNode(nnkEmpty),
    newNimNode(nnkEmpty),
    newNimNode(nnkEmpty),
    newNimNode(nnkFormalParams).add(newNimNode(nnkEmpty)),
    newNimNode(nnkEmpty),
    newNimNode(nnkEmpty),
    newNimNode(nnkStmtList).add(callExpr)
  ))

proc panic*(err: EmptyInterface) =
  raise (ref PanicException)(wrappedValue: err)

proc panic*(err: string) =
  let exc = newException(PanicException, err)
  exc.wrappedValue = err
  raise exc

proc recover*(): EmptyInterface =
  let err = getCurrentException()
  if err == nil:
    return null
  elif err of PanicException:
    exitDefer = true
    return (ref PanicException)(err).wrappedValue
  else:
    exitDefer = true
    return err

template godefer*(body: expr): stmt =
  let call = goPartialCall(body)
  defer:
    exitDefer = false
    call()
    if exitDefer:
      return
