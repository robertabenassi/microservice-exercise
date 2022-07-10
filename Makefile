fmt:                ## Format Go code
	@go fmt ./...

test:               ## Testing
	@go test ./internal/...

build:              ## builds portAPI and portDomainService
	@docker-compose build --no-cache

up:                 ## creates all containers needed by the environment
	@docker-compose up --force-recreate