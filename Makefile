-include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/app
run/app:
	@go run ./cmd/app -discord-token=${DISCORD_TOKEN} $(FLAGS) 
	
# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/app: build the cmd/app application
.PHONY: build/app
build/app:
	@go build -o bin/app ./cmd/app
