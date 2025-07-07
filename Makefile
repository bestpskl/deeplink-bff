GO := GO111MODULE=on go
DOCKER := DOCKER_DEFAULT_PLATFORM=linux/amd64

build:
	$(GO) build -mod=vendor -a -installsuffix cgo -tags musl -o main ./bff/cmd/deeplink-api/main.go

run-api:
	go run bff/cmd/deeplink-api/main.go

gen-swag:
	swag init --pd -d bff/cmd/deeplink-api -o ./docs 
	swag fmt -g bff/cmd/deeplink-api/main.go

clean:
	@rm -rf main ./vendor
