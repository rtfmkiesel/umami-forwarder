.PHONY: make
make:
	GO_ENABLED=0 GOOS=linux \
	go build \
		-ldflags="-s -w \
			-X main.version=@`git describe --tags`" \
		-o ./umami-forwarder ./main.go

.PHONY: lint
lint:
	deadcode ./...
	modernize ./...
	goimports-reviser -format ./...
	golangci-lint run

.PHONY: updatepackages
updatepackages:
	go mod tidy
	go get ${shell go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all}

.PHONY: docker-image
docker-image:
	docker build . --tag ghcr.io/rtfmkiesel/umami-forwarder:${shell git describe --tags}
	docker build . --tag ghcr.io/rtfmkiesel/umami-forwarder:latest

.PHONY: docker-push
docker-push:
	docker push ghcr.io/rtfmkiesel/umami-forwarder:${shell git describe --tags}
	docker push ghcr.io/rtfmkiesel/umami-forwarder:latest