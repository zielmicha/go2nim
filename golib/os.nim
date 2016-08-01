include gosupport
import golib/io

exttypes:
  type
    File = struct(())

var stdin*: gcptr[File]
var stdout*: gcptr[File]
var stderr*: gcptr[File]

proc write*(f: gcptr[File], p: GoSlice[byte]): tuple[n: int, err: Error] =
  nil
