
var exitDefer {.threadvar.}: bool

type
  PanicException* = object of Exception
    wrappedValue*: EmptyInterface

macro goPartialCall(s: untyped): expr =
  # FIXME: evaluate arguments before
  #s.treeRepr.echo

  result = newNimNode(nnkLambda).add(
    newNimNode(nnkEmpty),
    newNimNode(nnkEmpty),
    newNimNode(nnkEmpty),
    newNimNode(nnkFormalParams).add(newNimNode(nnkEmpty)),
    newNimNode(nnkEmpty),
    newNimNode(nnkEmpty),
    newNimNode(nnkStmtList).add(s)
  )

proc panic*(err: EmptyInterface) =
  raise (ref PanicException)(wrappedValue: err)

proc panic*(err: string) =
  let val: EmptyInterface = err
  panic(val)

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
