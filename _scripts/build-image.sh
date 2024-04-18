#!/bin/bash

# This script builds and optionally pushes a Docker image for multiple architectures using docker buildx.
# It's designed to be used in environments where Dockerfiles are located in various directories.

# Exit immediately if any command fails.
set -e

# Define default environment variables.
: "${PUSH="false"}"
: "${DOCKER_REPOSITORY="kperreau/goac"}"
: "${PROJECT_PATH="."}"  # Directory containing the Dockerfile.

# Define the Dockerfile path.
DOCKERFILE="${PROJECT_PATH}/Dockerfile"

# Get the current git commit hash shortened to 7 characters.
GIT_VERSION=$(git rev-parse --short=7 HEAD)

# Start constructing the docker build command.
buildCmd=(docker buildx build --platform="linux/amd64,linux/arm64" --network host)

# Append the --push option if the image should be pushed to a repository.
if [[ "${PUSH}" == "true" ]]; then
    buildCmd+=(--push)
fi

# Echo the full docker command for verification before execution.
echo "Executing Docker command:"
echo "${buildCmd[@]}" -t "${DOCKER_REPOSITORY}:latest" -t "${DOCKER_REPOSITORY}:${GIT_VERSION}" -f "${DOCKERFILE}" .

# Execute the docker build command.
"${buildCmd[@]}" -t "${DOCKER_REPOSITORY}:latest" -t "${DOCKER_REPOSITORY}:${GIT_VERSION}" -f "${DOCKERFILE}" .