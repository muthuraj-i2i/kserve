#!/bin/bash

# Copyright 2022 The KServe Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The script is used to build all the KServe images.

# TODO: Implement selective building and tag replacement based on modified code.

#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

echo "Github SHA ${GITHUB_SHA}"

# Define image tags
CONTROLLER_IMG_TAG=${DOCKER_REPO}/${CONTROLLER_IMG}:${GITHUB_SHA}
LOCALMODEL_CONTROLLER_IMG_TAG=${DOCKER_REPO}/${LOCALMODEL_CONTROLLER_IMG}:${GITHUB_SHA}
LOCALMODEL_AGENT_IMG_TAG=${DOCKER_REPO}/${LOCALMODEL_AGENT_IMG}:${GITHUB_SHA}
STORAGE_INIT_IMG_TAG=${DOCKER_REPO}/${STORAGE_INIT_IMG}:${GITHUB_SHA}
AGENT_IMG_TAG=${DOCKER_REPO}/${AGENT_IMG}:${GITHUB_SHA}
ROUTER_IMG_TAG=${DOCKER_REPO}/${ROUTER_IMG}:${GITHUB_SHA}

# Function to build & export as OCI tarball
build_and_export() {
  name=$1
  dockerfile=$2
  context_dir=$3

  echo "Building ${name} image"
  docker buildx build -f "${dockerfile}" "${context_dir}" \
    -t "${DOCKER_REPO}/${name}:${GITHUB_SHA}" \
    -o "type=oci,dest=${DOCKER_IMAGES_PATH}/${name}-${GITHUB_SHA}.tar" &
}

# Build all images in parallel
build_and_export "${CONTROLLER_IMG}" Dockerfile .
build_and_export "${LOCALMODEL_CONTROLLER_IMG}" localmodel.Dockerfile .
build_and_export "${LOCALMODEL_AGENT_IMG}" localmodel-agent.Dockerfile .
build_and_export "${AGENT_IMG}" agent.Dockerfile .
build_and_export "${ROUTER_IMG}" router.Dockerfile .
build_and_export "${STORAGE_INIT_IMG}" storage-initializer.Dockerfile python

wait

echo "All images built and exported as OCI tarballs"