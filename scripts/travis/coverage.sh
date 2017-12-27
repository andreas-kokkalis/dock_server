#!/usr/bin/env bash

set -ex

echo "" > coverage.txt
for dir in $(go list ./... | grep -v /spec); do
	go test -race -coverprofile=profile.out -covermode=atomic $dir
	if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
	fi
done
