# justfile for greenhead

# Show this justfile's commands.
list:
	just --list

# Prepare assets sources.  Requires binsanity.
assets:
	cat .github/rm-head.md > .github/README.md
	cat README.md | sed -n '/^<!-- cut -->$/,$p' | sed '1d' >> .github/README.md
	cat .github/rm-foot.md >> .github/README.md
	grep -v '^!\[' README.md > ghd/assets/src/doc/readme.md
	cd ghd/assets && binsanity src

# Run locally from source, with args passed. Args must not contain spaces.
run *ARGS='--version':
	go run ./ghd/cmd/ghd {{ARGS}}

# Serve the API locally, using a working test config; open a browser on Mac.
serve *ARGS:
	(which open && sleep 2 && open http://localhost:3030) &
	go run ./ghd/cmd/ghd api serve --config=testdata/config-full.toml {{ARGS}}

# As serve, but calls webui and assets to rebuild the SPA first.
serve-webui *ARGS: webui assets
	go run ./ghd/cmd/ghd api serve --config=testdata/config-full.toml {{ARGS}}

# Build for current environment.
build:
	go build -o build/ghd ./ghd/cmd/ghd

# Build for Apple Silicon Mac.
build-macos-arm:
	mkdir -p build/darwin-arm64
	GOOS=darwin GOARCH=arm64 go build -o build/darwin-arm64/ghd ./ghd/cmd/ghd

# Build for Intel Mac.
build-macos-intel:
	mkdir -p build/darwin-amd64
	GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/ghd ./ghd/cmd/ghd

# Build for Mac (all architectures).
build-macos: build-macos-arm build-macos-intel

# Build for Windows (all architectures).
build-windows: build-windows-arm build-windows-intel

# Build for Linux (all architectures).
build-linux: build-linux-arm build-linux-intel

# Build for all supported operating systems and architectures.
build-all: build-macos build-windows build-linux

# Build for ARM Windows.
build-windows-arm:
	mkdir -p build/windows-arm64
	GOOS=windows GOARCH=arm64 go build -o build/windows-arm64/ghd.exe

# Build for Intel Windows.
build-windows-intel:
	mkdir -p build/windows-amd64
	GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/ghd.exe

# Build for ARM Linux.
build-linux-arm:
	mkdir -p build/linux-arm64
	GOOS=linux GOARCH=arm64 go build -o build/linux-arm64/ghd.exe

# Build for Intel Linux.
build-linux-intel:
	mkdir -p build/linux-amd64
	GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/ghd.exe

# Make a full build with testing etc for all targets, CI-style.
build-full: clean vet license tooldoc test build-all
	@echo TODO: publish the builds somewhere
	@echo TODO: run coverage and upload it somewhere

# Release using goreleaser (will build everything to clean dist dir).
[confirm("⚠️ BUILD AND RELEASE ARTIFACTS TO GITHUB?")]
release:
	goreleaser release --clean

# Check for programmer errors.
vet:
	go vet ./...

# Run unit tests.
test:
	go test ./...

# Run unit tests with fresh assets.
atest: assets
    go test ./...

# Run unit tests with coverage, and open the coverage report.
cover:
	go test ./... -coverprofile=cover.out && go tool cover -html=cover.out

# Run tests with coverage for upload to Codecov.
codecov:
	go test ./... -coverprofile=coverage.txt

# Run benchmarks, if any.
bench:
	go test ./... -bench=.

# Run pkgsite because godoc is deprecated. :-(
doc:
	pkgsite -open

# Remove the build, release (dist) and cover artifacts.
clean:
	/bin/rm -rf build dist cover.out

# Generate the licenses for the doc command, checking compatibility first.
license:
	go-licenses check ./...
	mkdir -p build
	rm -rf build/third_party_licenses
	go-licenses save ./... --save_path=build/third_party_licenses
	misc/licenses.sh build/third_party_licenses ghd/assets/src/doc/licenses.md
	just assets

# Compile and bundle the webui files, putting them under assets.
webui:
	cd webui && prettier -w index.html src/*.*
	cd webui && npm run build
	cp -f webui/dist/index.html ghd/assets/src/webui/app.html
	just assets

# Copy README.md files from tools into assets.
tooldoc:
	find tools -name README.md -exec sh -c \
	'for f; do dir=$(dirname "$f"); base=$(basename "$dir"); \
	out="assets/src/doc/$dir.md"; mkdir -p "$(dirname "$out")"; \
	cp "$f" "$out"; done' \
	sh {} +
	just assets
