all: angeltrax

.PHONY: angeltrax
angeltrax: build/angeltrax build/angeltrax.exe

build:
	mkdir -p build

.PHONY: clean
clean:
	rm -rf build

.PHONY: test
test:
	go vet ./...
	go test ./...

build/angeltrax: build
	CGO_ENABLED=0 GOOS=linux go build -o $@ ./cmd/angeltrax/*.go

build/angeltrax.exe: build
	CGO_ENABLED=0 GOOS=windows go build -o $@ ./cmd/angeltrax/*.go
