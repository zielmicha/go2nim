#!/bin/bash
mkdir -p build/tests
cat gotests.txt | while read name; do
    mkdir -p "build/tests/$(dirname "$name.out")"
    ./gorun.sh go/test/$name >"build/tests/$name.out" 2>&1
    if [ $? = 0 ]; then
        rm "build/tests/$name.out"
        echo "$name ok"
    else
        echo "$name failed"
    fi
done
