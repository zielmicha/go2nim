#!/bin/sh
set -e

./gobuild.sh
./genstdlib.sh
echo 'import golib/strings, golib/fmt' | nim c -
