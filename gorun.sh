#!/bin/sh
set -e
./go2nim all "$@" > gorun.nim
nim c gorun.nim
echo "--- Running ---"
./gorun
echo "--- Exit code: $? ---"
