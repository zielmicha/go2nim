#!/bin/bash
mkdir -p build/tests
cat gotests.txt | while read name; do
    echo "=====" "$name"
    mkdir -p "build/tests/$(dirname "$name.out")"
    ./gorun.sh go/test/$name >"build/tests/$name.out" 2>&1
    if [ $? = 0 ]; then
        rm "build/tests/$name.out"
        echo "ok"
    else
        echo "failed"
    fi
done
