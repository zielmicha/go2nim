{.experimental.}
{.overflowChecks: off.}
# This file is included on top of all translated Go files.
import golib/compat, golib/builtin
import collections/gcptrs, collections/interfaces, collections/anonfield, collections/exttype, collections/goslice

# restore overflow checking for division
goSupportMakeMod SomeUnsignedInt
goSupportMakeMod int64
goSupportMakeMod int32
goSupportMakeMod int16
goSupportMakeMod int8
goSupportmakeMod int
