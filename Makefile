.PHONY: bootstrap
bootstrap:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u github.com/kyoh86/richgo
	go get -v github.com/ramya-rao-a/go-outline
	go get -v github.com/mdempsky/gocode
	go get -v github.com/uudashr/gopkgs/cmd/gopkgs
	go get -v golang.org/x/tools/cmd/goimports
	go get -v golang.org/x/tools/cmd/goimports

.PHONY: lint
lint: fmt
	golangci-lint run

.PHONY: test-flake
test-ci:
	@echo "tests:"
	richgo test -count=30 -v -cover ./...

.PHONY: install-golang-ci
lint-ci:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh

.PHONY: generate
generate:
	protoc --proto_path=orcareduce/model --go_out=orcareduce/model --go_opt=paths=source_relative datamodel.proto

.PHONY: fmt
fmt:
	@echo "fmt:"
	scripts/fmt

