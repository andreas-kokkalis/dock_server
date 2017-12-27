#!/usr/bin/env bash

set -e

echo "" > coverage.txt
for dir in $(go list ./... | grep -v /spec); do
	go test -v -coverprofile=profile.out -covermode=atomic $dir ; test ${PIPESTATUS[0]} -eq 0
	if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
	fi
done
