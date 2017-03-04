TOPDIR:=$(shell pwd)

#############
# VENDORING #
#############
# It removes vendor directories downloaded by glide
fix-vendor: rm -rf vendor/github.com/docker/docker/vendor

########################
# INITIALIZING PROJECT #
########################
bootstrap:
	cd scripts/seed_image
	make docker-build
	make docker-name
	cd ${TOPDIR}

##############
# TEST SUITE #
##############
# Go commands
GO_CMD=go
GO_TEST=$(GO_CMD) test -cover
# Packages to be tested
GO_TEST_PKG_LIST:= conf

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
	@go test -v --cover $(shell go list ./... | grep -v /vendor/) | sed ''/PASS/s//$(shell printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(shell printf "\033[31mFAIL\033[0m")/''
	make post-test
