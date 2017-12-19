#!/usr/bin/env bash

set -e

# Load function that checks if postgres is up and accepting connections
source scripts/docker/postgres_check.sh

set -x

docker-compose stop && yes | docker-compose rm
docker-compose up -d postgres redis

is_up=$(postgres_is_up)
if [[ $is_up == 'false' ]]; then
    echo "Postgres is not accepting connections yet."
    exit 1
fi

go test -v ./pkg/api/image/spec -ginkgo.v
