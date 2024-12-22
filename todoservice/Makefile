IMAGE_NAME ?= todoservice
TAG ?= latest
REGISTRY ?= victoruzunov

# Full image name
FULL_IMAGE_NAME := $(REGISTRY)/$(IMAGE_NAME):$(TAG)

# Default target
.PHONY: all
all: build push

# Build the Docker image
.PHONY: build
build:
	@echo "Building Docker image: $(FULL_IMAGE_NAME)"
	docker build -t $(FULL_IMAGE_NAME) .

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