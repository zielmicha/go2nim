# Translational module

proc println(args: varargs[string, `$`]) =
  echo(args)

proc print(args: varargs[string, `$`]) =
  echo(args)
