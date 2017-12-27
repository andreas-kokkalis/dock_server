TOPDIR:=$(shell pwd)

#############
# VENDORING #
#############
vendor:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

##############
# Unit Tests #
##############
unit-tests:
	@./scripts/travis/tests.sh

#####################
# Integration Tests #
#####################
integration-run:
	@./scripts/travis/integration-run.sh

#########################
# Update mock functions #
#########################
mock:
	go get github.com/matryer/moq
	go generate pkg/drivers/redis/redis.go

	mkdir -p pkg/api/repositories/redis && cp pkg/api/repositories/redis_repo.go pkg/api/repositories/redis/ && go generate pkg/api/repositories/redis/redis_repo.go && rm -rf pkg/api/repositories/redis/
	mkdir -p pkg/api/repositories/admin && cp pkg/api/repositories/db_admin_repo.go pkg/api/repositories/admin/ && go generate pkg/api/repositories/admin/db_admin_repo.go && rm -rf pkg/api/repositories/admin/
	# mkdir -p pkg/api/repositories/docker && cp pkg/api/repositories/docker_repo.go pkg/api/repositories/docker/ && go generate pkg/api/repositories/docker/docker_repo.go && rm -rf pkg/api/repositories/docker/
