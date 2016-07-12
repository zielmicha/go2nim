## Fake synchronization primitives.
import collections/exttype, collections/gcptrs

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
