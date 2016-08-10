include gosupport
import golib/io
import golib/errors

exttypes:
  type
    File = struct((
      sysfile: system.File
    ))

var stdin*: gcptr[File] = (ref File)(sysfile: system.stdin)
var stdout*: gcptr[File] = (ref File)(sysfile: system.stdout)
var stderr*: gcptr[File] =  (ref File)(sysfile: system.stderr)

proc write*(f: gcptr[File], p: GoSlice[byte]): tuple[n: int, err: Error] =
  let len = f.sysfile.writeBuffer(addr p[0], p.len)
  # TODO: handle errors
  result.n = len

proc read*(f: gcptr[File], p: GoSlice[byte]): tuple[n: int, err: Error] =
  panic("read")
