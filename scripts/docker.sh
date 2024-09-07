#!/usr/bin/env bash

HUB=${HUB:-"adhp"}
IMG=${IMG:-"bot"}
TAG=${TAG:-"dev"}
BASE_DISTRIBUTION=${BASE_DISTRIBUTION:-"distroless"}

docker buildx build \
  --platform linux/amd64 \
  --build-arg BASE_DISTRIBUTION=${BASE_DISTRIBUTION} \
  -t ${HUB}/${IMG}:${TAG} \
  -f Dockerfile \
  .