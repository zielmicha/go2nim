## Fake synchronization primitives.
import collections/exttype, collections/gcptrs
import golib/builtin

exttypes:
  type
    Mutex* = struct(())

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

# Pool
# TODO: thread safety

proc get*(p: gcptr[Pool]): EmptyInterface {.gomethod.} =
  if p.values.len == 0:
    p.values.add((p.new)())
  return p.values.pop()

proc put*(p: gcptr[Pool], x: EmptyInterface): void {.gomethod.} =
  if p.values.len < 32:
    p.values.add(x)
