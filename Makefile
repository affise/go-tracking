# suppress output, run `make XXX V=` to be verbose
V := @

default: test

.PHONY: clean
clean:
	$(V)golangci-lint cache clean

.PHONY: test
test: GO_TEST_FLAGS += -race
test:
	$(V)go test $(GO_TEST_FLAGS) --tags=$(GO_TEST_TAGS)  ./...

.PHONY: fulltest
fulltest: GO_TEST_TAGS += integration
fulltest: test

.PHONY: bench
bench: GO_TEST_FLAGS += -bench=. -benchmem -run=XXX -benchtime=20s
bench: test

.PHONY: lint
lint:
	$(V)golangci-lint run --config configs/.golangci.yml
	$(V)prototool lint $(PROTOTOOL_FLAGS)

.PHONY: generate
generate:
	$(V)prototool generate $(PROTOTOOL_FLAGS)
	$(V)go generate -x ./...

.PHONY: vendor
vendor:
	$(V)go mod tidy
	$(V)go mod vendor
	$(V)git add vendor

.PHONY: docker-lint
docker-lint:
	$(V)$(call in_gobuilder,make lint GOFLAGS=-buildvcs=false PROTOTOOL_FLAGS=--walk-timeout=15s)

.PHONY: docker-generate
docker-generate:
	$(V)$(call in_gobuilder,make generate PROTOTOOL_FLAGS=--walk-timeout=150s)

CURR_REPO := /$(notdir $(PWD))
define in_gobuilder
	docker run --rm \
		-v $(PWD):$(CURR_REPO) \
		-e CGO_ENABLED=0 \
		-w $(CURR_REPO) \
		docker.affisecorp.com/library/gobuilder:v1.21.0 $1
endef
