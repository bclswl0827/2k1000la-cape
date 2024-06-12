.PHONY: build clean run

BINARY=monitor
SRC_DIR=.
DIST_DIR=./build/dist
ASSETS_DIR=./build/assets

BUILD_ARCH=arm arm64 386 amd64 ppc64le riscv64 s390x loong64
BUILD_FLAGS=-s -w
BUILD_ARGS=-trimpath

build: clean $(BUILD_ARCH)
$(BUILD_ARCH):
	@echo "Building Linux $@ ..."
	@mkdir -p $(DIST_DIR)/$@
	@rm -rf $(DIST_DIR)/$@/*
	@CGO_ENABLED=0 GOOS=linux GOARCH=$@ go build -ldflags="$(BUILD_FLAGS)" \
		$(BUILD_ARGS) -o $(DIST_DIR)/$@/$(BINARY) $(SRC_DIR)/*.go
	@cp -r $(ASSETS_DIR) $(DIST_DIR)/$@

run:
	@go run $(SRC_DIR)/*.go --config $(ASSETS_DIR)/config.json

clean:
	@rm -rf $(DIST_DIR)/*
