#####################################
#
# GO 
#
#####################################
fmt:                ## Format Go code
	@go fmt ./...

test:               ## Testing
	@go test ./internal/...

#####################################
#
# DOCKER COMPOSE
#
#####################################
build:              ## builds portAPI and portDomainService
	@docker-compose build --no-cache

up:                 ## creates all containers needed by the environment
	@docker-compose up --force-recreate

#####################################
#
# MICROSERVICE  UPLOAD 
#
#####################################
upload:         ## Usage: make load-ports file=myports.json, by default it reads the testdata/ports.json file when a file option is not given
	@if [ "$(file)" = "" ]; then \
		curl -F file=@testdata/ports.json 'http://127.0.0.1:8000/updatePorts'; \
	else \
		curl -F file=@$(file) 'http://127.0.0.1:8000/updatePorts'; \
	fi;