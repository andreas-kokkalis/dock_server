#!/usr/bin/env bash

set -e

go get golang.org/x/tools/cmd/cover
go get github.com/axw/gocov/gocov
go get github.com/modocache/gover
go get github.com/mattn/goveralls

for dir in $(go list ./... | grep -v /vendor/ | grep -v /spec);
do
    go test -coverprofile=profile.coverprofile $dir
    echo "$test : status: $?"
done
gover
goveralls -coverprofile=gover.coverprofile -service=travis-ci
