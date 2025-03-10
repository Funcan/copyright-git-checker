default: audit test build


.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

.PHONY: test
test:
	go test -v -race -buildvcs -coverprofile coverage.out  ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover: test
	go tool cover -html=coverage.out

.PHONY: build
build:
	goreleaser build --snapshot --clean

.PHONY: clean
clean:
	rm -rf dist coverage.out
