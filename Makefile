BINARY_NAME=moose
BUILD_DIR=build
MODULE_PATH=./cmd/moose/

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-s -w -X main.Version=$(VERSION) -buildid=

.PHONY: build release run clean

build:
	@mkdir -p $(BUILD_DIR)
	go build -v -gcflags="all=-N -l" -ldflags="$(LDFLAGS)" -o "$(BUILD_DIR)/$(BINARY_NAME)-debug" $(MODULE_PATH)

release:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="$(LDFLAGS)" -o "$(BUILD_DIR)/$(BINARY_NAME)" $(MODULE_PATH)

run:
	go run $(MODULE_PATH) 

clean:
	rm -rf $(BUILD_DIR)
