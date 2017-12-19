#!/usr/bin/env bash

set -ex

: ${DOCKER_USER?"Set the environment variable DOCKER_USER"}
: ${DOCKER_PWD?"Set the environment variable DOCKER_PWD"}

docker login -u ${DOCKER_USER} -p ${DOCKER_PWD}
docker pull andreaskokkalis/dc --all-tags
