# Project variables
BINARY_NAME=system-services-ui
VERSION=0.1.0
DIST_NAME=$(BINARY_NAME)-$(VERSION)

# Go build flags
GO_FLAGS=-v

.PHONY: all build clean test archive install

all: build

## build: Compile the binary
build:
	go build $(GO_FLAGS) -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

## clean: Remove build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f $(DIST_NAME).tar.gz
	rm -rf dist/

## test: Run go tests
test:
	go test ./...

## archive: Create the source tarball for Fedora packaging
archive:
	# Create a temporary directory to structure the tarball correctly
	mkdir -p dist/$(DIST_NAME)
	# Copy all tracked git files to the dist directory
	git archive --format=tar HEAD | tar -x -C dist/$(DIST_NAME)
	# Create the tar.gz
	tar -czf $(DIST_NAME).tar.gz -C dist $(DIST_NAME)
	# Clean up temp dir
	rm -rf dist/
	@echo "Archive created: $(DIST_NAME).tar.gz"

## install: Install binary (useful for local testing outside of RPM)
install: build
	install -Dm0755 $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)