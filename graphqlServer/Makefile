IMAGE_NAME ?= graphql
TAG ?= latest
REGISTRY ?= victoruzunov

# Full image name
FULL_IMAGE_NAME := $(REGISTRY)/$(IMAGE_NAME):$(TAG)

# Default target
.PHONY: all
all: build push

ifndef GITHUB_TOKEN
  $(error GITHUB_TOKEN is not set. Please set it in your environment.)
endif

# Build the Docker image
.PHONY: build
build:
	@echo "Building Docker image: $(FULL_IMAGE_NAME)"
	docker build -t $(FULL_IMAGE_NAME) --build-arg GIT_PRIVATEMODULESGITTOKEN=$(GITHUB_TOKEN) .

# Tag the Docker image (optional step if you need to retag)
.PHONY: tag
tag:
	@echo "Tagging Docker image as $(FULL_IMAGE_NAME)"
	docker tag $(IMAGE_NAME):$(TAG) $(FULL_IMAGE_NAME)

# Push the Docker image
.PHONY: push
push: build
	@echo "Pushing Docker image: $(FULL_IMAGE_NAME)"
	docker push $(FULL_IMAGE_NAME)