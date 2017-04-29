#!/usr/bin/env bash

set -e

go get golang.org/x/tools/cmd/cover
go get github.com/mattn/goveralls

for dir in $(go list ./... | grep -v /vendor/ | grep -v /spec);
do
    go test -covermode=count -coverprofile=profile.tmp $dir
    echo "$test : status: $?"
    if [ -f profile.tmp ]
    then
        cat profile.tmp | tail -n +2 >> profile.out
        rm profile.tmp
    fi
done
goveralls -coverprofile=profile.out -service=travis-ci

