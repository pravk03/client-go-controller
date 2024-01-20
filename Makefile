
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: build
build: fmt vet vendor
	go build -o bin/main main.go

.PHONY: run
run:
	go run ./main.go