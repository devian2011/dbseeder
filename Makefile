
PHONY: test
test:
	go test ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: docker-build
	docker build -f ./docker/Dockerfile .  -t devian/dbseeder

