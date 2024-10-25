include .env

deploy:
	TAG=${TAG} && docker compose up -d
dev:
	docker build . -t hmdockerhub/price-api:dev --build-arg GITLAB_USERNAME=${GITLAB_USERNAME} --build-arg GITLAB_PASSWORD=${GITLAB_PASSWORD} --build-arg GO_MOD_TAG=${GO_MOD_TAG} --target=price-api
	docker compose -f docker-compose.dev.yaml up -d
clean:
	rm -rf build/
	go clean .
dev-routine:
	wire github.com/ipanardian/price-api/internal/injector
prod:
	make build-linux
get-module:
	go mod tidy
build-linux:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/api cmd/api/main.go
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/consumer cmd/consumer/main.go
