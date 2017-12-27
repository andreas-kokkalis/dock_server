[![Build Status](https://travis-ci.org/andreas-kokkalis/dock_server.svg?branch=master)](https://travis-ci.org/andreas-kokkalis/dock_server/)
[![codecov](https://codecov.io/gh/andreas-kokkalis/dock_server/branch/master/graph/badge.svg)](https://codecov.io/gh/andreas-kokkalis/dock_server)
[![Go Report Card](https://goreportcard.com/badge/github.com/andreas-kokkalis/dock_server)](https://goreportcard.com/report/github.com/andreas-kokkalis/dock_server)
# Dock-server
Dock server is an LTI tool provider, that allows creating docker images, running containers.

## Architecture
It relies on Docker Remote API, to communicate with the docker daemon, and manage the lifecycle of containers, and also manipulation of images.

It has a RESTful api, that exposes such functinality. The API is consumed by a client application, named dock-client, which provides the LTI integrations, as long as frontend for the admin of dock-server.

### Software Dependencies
Docker version >=1.12.3

glide version 0.13.0-dev https://github.com/Masterminds/glide

golang 1.7
