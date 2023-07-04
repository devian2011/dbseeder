
PHONY: test
test:
	go test ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: docker-build
	docker build -f ./docker/Dockerfile .  -t devian/dbseeder

.PHONY: lint
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.52.2 golangci-lint run -v
