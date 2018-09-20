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

IMAGE = $(REGISTRY)/openshift-state-metrics
MULTI_ARCH_IMG = $(IMAGE)-$(ARCH)

gofmtcheck:
	@go fmt $(PKGS) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi
build: clean
	docker run --rm -v "$$PWD":/go/src/github.com/openshift/openshift-state-metrics -w /go/src/github.com/openshift/openshift-state-metrics -e GOOS=$(shell uname -s | tr A-Z a-z) -e GOARCH=$(ARCH) -e CGO_ENABLED=0 golang:${GO_VERSION} go build -ldflags "-s -w -X ${PKG}/version.Release=${TAG} -X ${PKG}/version.Commit=${Commit} -X ${PKG}/version.BuildDate=${BuildDate}" -o openshift-state-metrics

test-unit: clean build
	GOOS=$(shell uname -s | tr A-Z a-z) GOARCH=$(ARCH) $(TESTENVVAR) go test --race $(FLAGS) $(PKGS)

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
	docker run --rm -v "$$PWD":/go/src/github.com/openshift/openshift-state-metrics -w /go/src/github.com/openshift/openshift-state-metrics -e GOOS=linux -e GOARCH=$(ARCH) -e CGO_ENABLED=0 golang:${GO_VERSION} go build -ldflags "-s -w -X ${PKG}/version.Release=${TAG} -X ${PKG}/version.Commit=${Commit} -X ${PKG}/version.BuildDate=${BuildDate}" -o openshift-state-metrics
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

clean:
	rm -f openshift-state-metrics

.PHONY: all build all-push all-container test-unit container quay-push clean
