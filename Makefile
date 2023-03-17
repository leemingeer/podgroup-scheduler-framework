
COMMONENVVAR=GOOS=$(shell uname -s | tr A-Z a-z)
BUILDENVVAR=CGO_ENABLED=0


# version
RELEASE_VERSION?=v$(shell date +%Y%m%d)-$(shell git describe --tags --match "v*")
# The RELEASE_VERSION variable can have one of two formats:
# v20201009-v0.18.800-46-g939c1c0 - automated build for a commit(not a tag) and also a local build
# v20200521-v0.18.800             - automated build for a tag
VERSION=$(shell echo $(RELEASE_VERSION) | awk -F - '{print $$2}')

# image
LOCAL_REGISTRY=localhost:5000/scheduler-plugins
RELEASE_REGISTRY?=leemingeer
RELEASE_IMAGE:=sample-scheduler:$(RELEASE_VERSION)

# compile
.PHONY: build-scheduler.amd64
build-scheduler.amd64:
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags '-X k8s.io/component-base/version.gitVersion=$(VERSION) -w' -o bin/sample-scheduler cmd/scheduler/main.go

.PHONY: build-scheduler.arm64
build-scheduler.arm64:
	GOOS=linux CGO_ENABLED=0 GOARCH=arm64 go build -ldflags '-X k8s.io/component-base/version.gitVersion=$(VERSION) -w' -o bin/sample-scheduler cmd/scheduler/main.go


# image
.PHONY: local-image
local-image: clean
	docker build -f ./build/scheduler/Dockerfile --build-arg ARCH="arm64" --build-arg RELEASE_VERSION="$(RELEASE_VERSION)" -t $(LOCAL_REGISTRY)/$(LOCAL_IMAGE) .

.PHONY: release-image.amd64
release-image.amd64: clean
	docker build -f ./build/scheduler/Dockerfile --build-arg ARCH="amd64" --build-arg RELEASE_VERSION="$(RELEASE_VERSION)" -t $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-amd64 .

.PHONY: release-image.arm64
release-image.arm64: clean
	docker build -f ./build/scheduler/Dockerfile --build-arg ARCH="arm64" --build-arg RELEASE_VERSION="$(RELEASE_VERSION)" -t $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-arm64 .

.PHONY: clean
clean:
	rm -rf ./bin