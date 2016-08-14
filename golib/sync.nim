## Fake synchronization primitives.
import collections/exttype, collections/gcptrs, collections/reflect
import golib/builtin

exttypes:
  type
    Mutex* = struct(())
    RWMutex* = struct(())

    Cond* = struct((
      l: Locker
    ))

    Locker* = iface((
      lock(): void,
      unlock(): void
    ))

    Pool* = struct((
      new: (proc(): EmptyInterface),
      values: seq[EmptyInterface]
    ))

proc lock*(m: Mutex) =
  discard

proc unlock*(m: Mutex) =
  discard

proc broadcast*(c: Cond) =
  discard

proc wait*(c: Cond) =
  discard

proc signal*(c: Cond) =
  discard

proc lock*(m: RWMutex) =
  discard

proc unlock*(m: RWMutex) =
  discard

proc rLock*(m: RWMutex) =
  discard

proc rUnlock*(m: RWMutex) =
  discard

# Pool
# TODO: thread safety

proc get*(p: gcptr[Pool]): EmptyInterface {.gomethod.} =
  if p.values == nil:
    p.values = newSeq[EmptyInterface]()
  if p.values.len == 0:
    let f = (p.new)
    let val = f()
    p.values.add(val)
  result = p.values.pop()

proc put*(p: gcptr[Pool], x: EmptyInterface): void {.gomethod.} =
  if p.values == nil:
    p.values = newSeq[EmptyInterface]()
  if p.values.len < 32:
    p.values.add(x)
