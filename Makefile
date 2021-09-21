.PHONY: test cover godog up down 

test:
	@echo "running all tests"
	go vet ./...
	go test -v ./... -race

cover:
	@echo "print coverage info to stdout"
	go test -v ./internal/f3 ./pkg/accounts -race -coverprofile=coverage.out
	go tool cover -func=coverage.out
	rm coverage.out

godog:
	cd ./test/accounts; \
	godog -f progress; \
	cd ../..

up:
	docker-compose up

down:
	docker-compose down
