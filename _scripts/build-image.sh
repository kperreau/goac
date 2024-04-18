#!/bin/bash

set -e

: "${PUSH_IMAGE="false"}"

: "${REPOSITORY="kperreau/goac"}"

: "${PROJECT_PATH="."}" # project path (must be where the Dockerfile is)

DOCKERFILE="${PROJECT_PATH}/Dockerfile"

GIT_VERSION=$(git rev-parse --short=7 HEAD)

dockerCmd=(docker buildx build --platform="linux/amd64,linux/arm64" --network host)

if [[ "${PUSH_IMAGE}" == "true" ]]; then
    dockerCmd+=(--push);
fi

# print cmd
echo "${dockerCmd[@]}" \
  -t "${REPOSITORY}:latest" \
  -t "${REPOSITORY}:${GIT_VERSION}" \
  -f "${DOCKERFILE}" \
  .

# run docker build
"${dockerCmd[@]}" \
  -t "${REPOSITORY}:latest" \
  -t "${REPOSITORY}:${GIT_VERSION}" \
  -f "${DOCKERFILE}" \
  .