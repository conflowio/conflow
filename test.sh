#!/usr/bin/env bash

set -e
echo -n "" > coverage.txt

for dir in $(go list ./...); do
    go test -v -race -coverprofile=profile.out -covermode=atomic "$dir"
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
