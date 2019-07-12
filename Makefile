FLAGS =
TESTENVVAR =
REGISTRY = quay.io/redhat
TAG = $(shell git describe --abbrev=0)
PKGS = $(shell go list ./... | grep -v /vendor/)
ARCH ?= $(shell go env GOARCH)
BuildDate = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
Commit = $(shell git rev-parse --short HEAD)
ALL_ARCH = amd64 arm arm64 ppc64le s390x
PKG=github.com/openshift/openshift-state-metrics/pkg
GO_VERSION=1.11
JB_BINARY:=$(firstword $(subst :, ,$(GOPATH)))/bin/jb
JSONNET_VENDOR=jsonnet/jsonnetfile.lock.json jsonnet/vendor
JSONNET_SRC=$(shell find ./jsonnet -type f)
GOJSONTOYAML_BINARY:=$(firstword $(subst :, ,$(GOPATH)))/bin/gojsontoyaml

IMAGE = $(REGISTRY)/openshift-state-metrics
MULTI_ARCH_IMG = $(IMAGE)-$(ARCH)

validate-modules:
	@echo "- Verifying that the dependencies have expected content..."
	GO111MODULE=on go mod verify
	@echo "- Checking for any unused/missing packages in go.mod..."
	GO111MODULE=on go mod tidy
	@echo "- Checking for unused packages in vendor..."
	GO111MODULE=on go mod vendor
	@git diff --exit-code -- go.sum go.mod vendor/

doccheck:
	@echo "- Checking if the documentation is up to date..."
	@grep -hoE '(openshift_[^ |]+)' docs/* --exclude=README.md| sort -u > documented_metrics
	@sed -n 's/.*# TYPE \(openshift_[^ ]\+\).*/\1/p' pkg/collectors/*_test.go | sort -u > tested_metrics
	@diff -u0 tested_metrics documented_metrics || (echo "ERROR: Metrics with - are present in tests but missing in documentation, metrics with + are documented but not tested."; exit 1)
	@echo OK
	@rm -f tested_metrics documented_metrics
	@echo "- Checking for orphan documentation files"
	@cd docs; for doc in *.md; do if [ "$$doc" != "README.md" ] && ! grep -q "$$doc" *.md; then echo "ERROR: No link to documentation file $${doc} detected"; exit 1; fi; done
	@echo OK

gofmtcheck:
	@go fmt $(PKGS) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi
build: clean
	GO111MODULE=on GOOS=$(shell uname -s | tr A-Z a-z) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "-s -w -X ${PKG}/version.Release=${TAG} -X ${PKG}/version.Commit=${Commit} -X ${PKG}/version.BuildDate=${BuildDate}" -o openshift-state-metrics
test-unit: clean build
	GO111MODULE=on GOOS=$(shell uname -s | tr A-Z a-z) GOARCH=$(ARCH) $(TESTENVVAR) go test --race $(FLAGS) $(PKGS)

TEMP_DIR := $(shell mktemp -d)

all: all-container

sub-container-%:
	$(MAKE) --no-print-directory ARCH=$* container

sub-push-%:
	$(MAKE) --no-print-directory ARCH=$* push

all-container: $(addprefix sub-container-,$(ALL_ARCH))

all-push: $(addprefix sub-push-,$(ALL_ARCH))

container: .container-$(ARCH)
.container-$(ARCH):
	docker run --rm -v "$$PWD":/go/src/github.com/openshift/openshift-state-metrics -w /go/src/github.com/openshift/openshift-state-metrics -e GO111MODULE=on  -e GOOS=linux -e GOARCH=$(ARCH) -e CGO_ENABLED=0 golang:${GO_VERSION} go build -ldflags "-s -w -X ${PKG}/version.Release=${TAG} -X ${PKG}/version.Commit=${Commit} -X ${PKG}/version.BuildDate=${BuildDate}" -o openshift-state-metrics
	cp -r * $(TEMP_DIR)
	docker build -t $(MULTI_ARCH_IMG):$(TAG) $(TEMP_DIR)
	docker tag $(MULTI_ARCH_IMG):$(TAG) $(MULTI_ARCH_IMG):latest

ifeq ($(ARCH), amd64)
	# Adding check for amd64
	docker tag $(MULTI_ARCH_IMG):$(TAG) $(IMAGE):$(TAG)
	docker tag $(MULTI_ARCH_IMG):$(TAG) $(IMAGE):latest
endif

quay-push: .quay-push-$(ARCH)
.quay-push-$(ARCH): .container-$(ARCH)
	docker push $(MULTI_ARCH_IMG):$(TAG)
	docker push $(MULTI_ARCH_IMG):latest
ifeq ($(ARCH), amd64)
	docker push $(IMAGE):$(TAG)
	docker push $(IMAGE):latest
endif

jsonnet/vendor: $(JB_BINARY) jsonnet/jsonnetfile.json
	cd jsonnet && jb install

manifests: $(JSONNET_SRC) $(JSONNET_VENDOR) $(GOJSONTOYAML_BINARY) hack/build-jsonnet.sh
	./hack/build-jsonnet.sh

$(JB_BINARY):
	go get -u github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb

$(GOJSONTOYAML_BINARY):
	go get -u github.com/brancz/gojsontoyaml

clean:
	rm -f openshift-state-metrics

.PHONY: all build all-push all-container test-unit container quay-push clean validate-modules
