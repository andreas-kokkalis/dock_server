TOPDIR:=$(shell pwd)

#############
# VENDORING #
#############
vendor:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

########################
# INITIALIZING PROJECT #
########################
# TODO: remove this, since integration tests pull the seed image from docker hub
# bootstrap:
# 	cd scripts/seed_image
# 	make docker-build
# 	make docker-name
# 	cd ${TOPDIR}


##############
# Unit Tests #
##############
unit-tests:
	@./scripts/travis/tests.sh

#####################
# Integration Tests #
#####################
# Perform unit tests of all packages
LOG_DIR := ${TOPDIR}/logs
pre-test:
	@if [ ! -d "${LOG_DIR}" ]; then \
		mkdir ${LOG_DIR} ;\
	fi
	@if [ -f "${LOG_DIR}/compose.log" ]; then \
		rm ${LOG_DIR}/compose.log ;\
	fi
	@touch ${LOG_DIR}/compose.log
	@yes | docker-compose rm
	export MODE=DEV
	docker-compose up -d
	export MODE=

post-test:
	docker-compose stop

tests: pre-test
	@go test -v --cover $(shell go list ./... ) | sed ''/PASS/s//$(shell printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(shell printf "\033[31mFAIL\033[0m")/''
	make post-test

#########################
# Update mock functions #
#########################
mock:
	go get github.com/matryer/moq
	go generate pkg/cache/redis.go
