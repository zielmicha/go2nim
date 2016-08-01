include gosupport
import golib/io
import golib/errors

exttypes:
  type
    File = struct(())

var stdin*: gcptr[File]
var stdout*: gcptr[File]
var stderr*: gcptr[File]

proc write*(f: gcptr[File], p: GoSlice[byte]): tuple[n: int, err: Error] =
  panic("write")

proc read*(f: gcptr[File], p: GoSlice[byte]): tuple[n: int, err: Error] =
  panic("read")
