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
    VERSION=1.0

    all: test build-linux

    build:
		$(GOBUILD) -o $(BINARY_NAME) -trimpath -ldflags "-X main.sha1ver=$(VERSIONGIT) -X main.buildTime=$(NOW)" -v

    test:
		$(GOTEST) -v ./...

    fmt:
		$(GOFMT) ./...

    clean:
		$(GOCLEAN)
		$(MODCLEAN)

    run:
		$(GOBUILD) -o build/ -trimpath -ldflags "-X main.sha1ver=$(VERSIONGIT) -X main.buildTime=$(NOW)" -v ./...
		./build/$(BINARY_NAME)

    deps:
		$(GOGET) github.com/marstid/go-ontap
		$(GOGET) github.com/spf13/viper

	# Depen
    deps-refresh:
		$(GOGET) -u github.com/marstid/go-ontap
		$(GOGET) -u github.com/spf13/viper

    # Cross compilation
    build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -trimpath -ldflags "-X main.sha1ver=$(VERSIONGIT) -X main.buildTime=$(NOW)" -v
