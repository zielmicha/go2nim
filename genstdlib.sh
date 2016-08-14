#!/bin/sh -e

genpkg() {
    echo $1
    ./go2nim all go/src/$1/*.go > stdlib/golib/$1.nim
}

mkdir -p stdlib/golib/unicode stdlib/golib/go

genpkg errors
genpkg io
genpkg unicode

for name in lo pc ps; do
    sed -i 's/  '$name'\* =/  `'$name'-`* =/' stdlib/golib/unicode.nim
    sed -i 's/: '$name',/: `'$name'-`,/' stdlib/golib/unicode.nim
done

genpkg unicode/utf8
genpkg strings
sed -i 's/if (var m = count(s, old); (m == 0)):/var m = count(s, old)\n  if (m == 0):/' stdlib/golib/strings.nim
genpkg strconv
genpkg fmt
genpkg sort
genpkg go/token
genpkg bytes

# FIXME: result can't be captured
sed -i 's/        result.err = se.err/        discard/' stdlib/golib/fmt.nim
sed -i 's/godefer(errorHandler(gcaddr result.err))//' stdlib/golib/fmt.nim
