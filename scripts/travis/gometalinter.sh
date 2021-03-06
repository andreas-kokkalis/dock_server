#!/usr/bin/env bash

set -ex

go get -u -v github.com/alecthomas/gometalinter
gometalinter --install

pkgs=$(go list -f "{{ .Dir }}" ./... | grep -v /vendor/ )
gometalinter --vendor --disable=gocyclo --disable=gas --dupl-threshold=70 --checkstyle --deadline=500s --json $pkgs | test ${PIPESTATUS[0]} -eq 0
