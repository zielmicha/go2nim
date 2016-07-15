#!/bin/sh -e

genpkg() {
    echo $1
    ./go2nim all ../go/src/$1/*.go > golib/$1.nim
    echo $1.nim >> golib/.gitignore
}

echo > golib/.gitignore
mkdir -p golib/unicode

genpkg errors
genpkg io
genpkg unicode
for name in lo pc ps; do
    sed -i 's/  '$name'\* =/  `'$name'-`* =/' golib/unicode.nim
    sed -i 's/: '$name',/: `'$name'-`,/' golib/unicode.nim
done
genpkg unicode/utf8
genpkg strings
genpkg strconv