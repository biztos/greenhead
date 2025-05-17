# justfile for greenhead

# Show this justfile's commands.
list:
	just --list

# Prepare assets etc.
prep:
	cp README.md assets/src/doc/readme.md && cp API.md assets/src/doc/api.md
	cd assets && binsanity src
	mkdir -p build

# Run locally from source, with args passed. Args must not contain spaces.
run *ARGS='--version': prep
	go run ./cmd/ghd {{ARGS}}

# Serve the API locally. Uses agent chatty by default.
serve *ARGS='--agent=chatty': prep
	go run ./cmd/ghd api serve {{ARGS}}

# As serve, but calls webui to rebuild the SPA first.
serve-webui *ARGS='--agent=chatty': webui prep
	go run ./cmd/ghd api serve {{ARGS}}

# Build for current environment.
build: prep
	go build -o build/ghd ./cmd/ghd

# Build for Apple Silicon Mac.
build-macos-arm: prep
	mkdir -p build/darwin-arm64
	GOOS=darwin GOARCH=arm64 go build -o build/darwin-arm64/ghd ./cmd/ghd

# Build for Intel Mac.
build-macos-intel: prep
	mkdir -p build/darwin-amd64
	GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/ghd ./cmd/ghd

# Build for Mac (all architectures).
build-macos: build-macos-arm build-macos-intel

# Build for Windows (all architectures).
build-windows: build-windows-arm build-windows-intel

# Build for Linux (all architectures).
build-linux: build-linux-arm build-linux-intel

# Build for all supported operating systems and architectures.
build-all: build-macos build-windows build-linux

# Build for ARM Windows.
build-windows-arm: prep
	mkdir -p build/windows-arm64
	GOOS=windows GOARCH=arm64 go build -o build/windows-arm64/ghd.exe

# Build for Intel Windows.
build-windows-intel: prep
	mkdir -p build/windows-amd64
	GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/ghd.exe

# Build for ARM Linux.
build-linux-arm: prep
	mkdir -p build/linux-arm64
	GOOS=linux GOARCH=arm64 go build -o build/linux-arm64/ghd.exe

# Build for Intel Linux.
build-linux-intel: prep
	mkdir -p build/linux-amd64
	GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/ghd.exe

# Make a full build with testing etc for all targets, CI-style.
build-full: clean vet license tooldoc test build-all
	@echo TODO: publish the builds somewhere
	@echo TODO: run coverage and upload it somewhere

# Check for programmer errors.
vet: prep
	go vet ./...

# Run unit tests.
test: prep
	go test ./...

# Run unit tests with coverage, and open the coverage report.
cover: prep
	go test ./... -coverprofile=cover.out && go tool cover -html=cover.out

# Run benchmarks, if any.
bench: prep
	go test ./... -bench=.

# Run pkgsite because godoc is deprecated. :-(
doc: prep
	pkgsite -open

# Remove the build and cover artifacts.
clean:
	/bin/rm -rf build cover.out

# Generate the licenses for the doc command, checking compatibility first.
license:
	go-licenses check ./...
	mkdir -p build
	rm -rf build/third_party_licenses
	go-licenses save ./... --save_path=build/third_party_licenses
	build-tools/licenses.sh build/third_party_licenses assets/src/doc/licenses.md

# Compile and bundle the webui files, putting them under assets.
webui:
	cd webui && prettier -w index.html src/*.*
	cd webui && npm run build
	cp -f webui/dist/index.html assets/src/webui/app.html

# Copy README.md files from tools into documenation source.
tooldoc:
	find tools -name README.md -exec sh -c \
	'for f; do dir=$(dirname "$f"); base=$(basename "$dir"); \
	out="assets/src/doc/$dir.md"; mkdir -p "$(dirname "$out")"; \
	cp "$f" "$out"; done' \
	sh {} +
