PROTOS			:= $(wildcard pkg/apis/*/*/*.proto)
ALL_SRC			:= $(shell find . -name "*.go" | grep -v -e vendor)
PACKAGES 		:= $(shell go list ./...)
PASS     		= $(shell printf "\033[32mPASS\033[0m")
FAIL     		= $(shell printf "\033[31mFAIL\033[0m")
COLORIZE 		= sed ''/PASS/s//$(PASS)/'' | sed ''/FAIL/s//$(FAIL)/''
HOSTNAME		=frankgreco
NAMESPACE		=ubiquiti
NAME			=edge
BINARY			=terraform-provider-${NAME}
VERSION			=0.0.1
OS_ARCH			=darwin_amd64
ACCTEST_TIMEOUT =5m
GO_OPTIONS		= -buildvcs=false

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: fmt
fmt:
	@gofmt -e -s -l -w $(ALL_SRC)

# Run 'brew install coreutils' to install sha256sum on Mac
define docs-generate-sum
	rm -f $@; \
	( \
		find templates internal/provider examples -name '*schema_*' -o -path '*examples*' -not -path '*.terraform*' -not -path 'examples/guides/*/provider.tf' -type f -o -path '*templates*' -type f | xargs sha256sum; \
	) | sort -k 2 > $@
endef

docs-generate.sum: docs-generate.sum.current
	@if cmp $@.current $@; then \
		echo "docs up-to-date"; \
	else \
		echo go generate ./...; \
		go generate ./...; \
		$(docs-generate-sum); \
	fi

docs-generate.sum.current: .FORCE
	@$(docs-generate-sum)

.PHONY: generate
generate: docs-generate.sum

.PHONY: deps
deps:
	@go mod download
	@go mod tidy

.PHONY: test
test: deps
	@bash -c "set -e; set -o pipefail; go test $(GO_OPTIONS) -v -race $(PACKAGES) | $(COLORIZE)"

testacc:
	TF_ACC=1 go test ./internal/provider/... -v -parallel 1 -timeout $(ACCTEST_TIMEOUT)

.PHONY: .FORCE
.FORCE:
