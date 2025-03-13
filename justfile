# justfile for greenhead

# Show this justfile's commands.
list:
	just --list

# Prepare assets etc.
prep:
	@echo prep is TBD, see if we have any assets first

# Run locally from source, with args passed.
run *ARGS='--version': prep
	go run ./cmd/ghd {{ARGS}}

# Build for current environment.
build: prep
	mkdir -p build
	go build -o build/ghd ./cmd/ghd

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

# Remove the build and cover artifacts.
clean:
	/bin/rm -rf build cover.out

