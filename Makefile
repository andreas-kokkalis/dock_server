# Go commands
GO_CMD=go
GO_TEST=$(GO_CMD) test -cover

# Packages
TOP_PACKAGE_DIR := ./
PACKAGE_LIST := conf


# It removes vendor directories downloaded by glide
fix-vendor: rm -rf vendor/github.com/docker/docker/vendor

# Perform unit tests of all packages
test:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Unit Testing $$p ..."; \
		$(GO_TEST) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done
