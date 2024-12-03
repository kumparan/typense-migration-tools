SHELL:=/bin/bash

ifdef test_run
	TEST_ARGS := -run $(test_run)
endif

test_command=richgo test ./... $(TEST_ARGS) -v --cover

test: lint test-only

check-gotest:
ifeq (, $(shell which richgo))
	$(warning "richgo is not installed, falling back to plain go test")
	$(eval TEST_BIN=go test)
else
	$(eval TEST_BIN=richgo test)
endif
ifdef test_run
	$(eval TEST_ARGS := -run $(test_run))
endif
	$(eval test_command=$(TEST_BIN) ./... $(TEST_ARGS) -v --cover)

test-only: check-gotest
	SVC_ENV=test SVC_DISABLE_CACHING=true $(test_command) -timeout 60s

check-cognitive-complexity:
	find . -type f -name '*.go' -exec gocognit -over 17 {} +

lint: check-cognitive-complexity
	golangci-lint run

changelog:
ifdef version
	$(eval changelog_args=--next-tag $(version) $(changelog_args))
	@echo $$(basename $$(git remote get-url origin) .git)@$(version) > VERSION
endif
	git-chglog $(changelog_args)

.PHONY: test lint test-only changelog
