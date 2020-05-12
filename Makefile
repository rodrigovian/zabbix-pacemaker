# Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    MODCLEAN=$(GOCMD) clean -modcache
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    GOFMT=$(GOCMD) fmt
    BINARY_NAME=zpmkr
    BINARY_LINUX=$(BINARY_NAME)_amd64
    NOW=`date +'%Y-%m-%d_%T'`
    VERSIONGIT=`git rev-parse HEAD`
    VERSION=`git describe --tags`

    all: test build-linux

    build:
		$(GOBUILD) -o $(BINARY_NAME) -trimpath -ldflags "-X main.sha1ver=$(VERSIONGIT) -X main.release=$(VERSION)  -X main.buildTime=$(NOW)" -v

    test:
		$(GOTEST) -v ./...

    fmt:
		$(GOFMT) ./...

    clean:
		$(GOCLEAN)
		$(MODCLEAN)
		rm -f $(BINARY_LINUX)
		rm -f $(BINARY_NAME)

    run:
		$(GOBUILD) -o build/ -trimpath -ldflags "-X main.sha1ver=$(VERSIONGIT) -X main.buildTime=$(NOW)" -v ./...
		./build/$(BINARY_NAME)

    # Cross compilation
    build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -trimpath -ldflags "-X main.sha1ver=$(VERSIONGIT) -X main.release=$(VERSION) -X main.buildTime=$(NOW)" -v
