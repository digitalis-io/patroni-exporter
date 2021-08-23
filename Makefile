VERSION ?= $(shell git tag | sort -V | tail -1)
APP_NAME=digitalisdocker/patroni-exporter
PORT=9877
DOCKER_USER=""
DOCKER_PASS=""
DOCKER_HOST=""

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help


# DOCKER TASKS
# Build the container
build: ## Build the container
	docker build -t $(APP_NAME) .

build-nc: ## Build the container without caching
	docker build --no-cache $(APP_NAME) .

run: ## Run container shell
	docker run -i -t --rm -v $(pwd):/config --entrypoint=/bin/bash --name="helm" $(APP_NAME)

up: build run ## Runs "build" and "run"

stop: ## Stop and remove a running container
	docker stop $(APP_NAME); docker rm $(APP_NAME)

tag: ## Tag contsiner with $VERSION
	docker tag $(APP_NAME) $(APP_NAME):$(VERSION)

login: ## Login to docker registry at $DOCKER_HOST using $DOCKER_USER / $DOCKER_PASS
	docker login --username=${DOCKER_USER} --password=${DOCKER_PASS} ${DOCKER_HOST}

push: tag login ## Push tag to registry
	docker push $(APP_NAME):$(VERSION)

version: ## Output the current version
	@echo $(VERSION)
	
