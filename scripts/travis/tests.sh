#!/usr/bin/env bash

set -ex

go test -v $(go list ./... | grep -v /vendor/ |grep -v /spec) ; test ${PIPESTATUS[0]} -eq 0
