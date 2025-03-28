.PHONY: build run deploy clean

# Variables
APP_NAME=whatstore-api
DOCKER_IMAGE=$(APP_NAME):latest

# Build the Docker image
build:
	docker build -t $(DOCKER_IMAGE) .

# Run the container
run:
	docker run -d --name $(APP_NAME) \
		-p 8080:8080 \
		--restart unless-stopped \
		$(DOCKER_IMAGE)

# Stop and remove the container
clean:
	-docker stop $(APP_NAME)
	-docker rm $(APP_NAME)

# Deploy - Rebuilds and runs the application
deploy: clean build run
	@echo "Deployment complete. API is running on port 8080"

# Show logs
logs:
	docker logs -f $(APP_NAME)
