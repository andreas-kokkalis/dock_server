#!/usr/bin/env bash

set -ex

docker login -u $DOCKER_HUB_USER -p $DOCKER_HUB_PWD
docker pull andreaskokkalis/dc --all-tags
