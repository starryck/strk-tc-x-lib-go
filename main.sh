#!/bin/bash
if [ -f .env ]; then
    export GIT_TAG="$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match)"
    export GIT_COMMIT="$(git log -1 --format=%h)"
fi

go build -o ./main.exe -tags="sonic avx" .
./main.exe "$@"
