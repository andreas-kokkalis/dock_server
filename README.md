[![Build Status](https://travis-ci.org/andreas-kokkalis/dock_server.svg?branch=master)](https://travis-ci.org/andreas-kokkalis/dock_server/)
[![Coverage Status](https://coveralls.io/repos/github/andreas-kokkalis/dock_server/badge.svg?branch=master)](https://coveralls.io/github/andreas-kokkalis/dock_server?branch=master)
# Dock-server
Dock server is an LTI tool provider, that allows creating docker images, running containers.

## Architecture
It relies on Docker Remote API, to communicate with the docker daemon, and manage the lifecycle of containers, and also manipulation of images.

It has a RESTful api, that exposes such functinality. The API is consumed by a client application, named dock-client, which provides the LTI integrations, as long as frontend for the admin of dock-server.

### Software Dependencies
Docker version >=1.12.3

glide version 0.13.0-dev https://github.com/Masterminds/glide

golang 1.7

