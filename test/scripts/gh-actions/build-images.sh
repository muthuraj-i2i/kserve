#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

echo "Github SHA ${GITHUB_SHA}"

# Image tags
CONTROLLER_IMG_TAG=${DOCKER_REPO}/${CONTROLLER_IMG}:${GITHUB_SHA}
LOCALMODEL_CONTROLLER_IMG_TAG=${DOCKER_REPO}/${LOCALMODEL_CONTROLLER_IMG}:${GITHUB_SHA}
LOCALMODEL_AGENT_IMG_TAG=${DOCKER_REPO}/${LOCALMODEL_AGENT_IMG}:${GITHUB_SHA}
STORAGE_INIT_IMG_TAG=${DOCKER_REPO}/${STORAGE_INIT_IMG}:${GITHUB_SHA}
AGENT_IMG_TAG=${DOCKER_REPO}/${AGENT_IMG}:${GITHUB_SHA}
ROUTER_IMG_TAG=${DOCKER_REPO}/${ROUTER_IMG}:${GITHUB_SHA}

# Cache directories (local cache, you can switch to registry cache in CI)
CACHE_DIR="/tmp/.buildx-cache"
CACHE_DIR_NEW="/tmp/.buildx-cache-new"

mkdir -p "${DOCKER_IMAGES_PATH}"

# Helper function to build with cache
build_image() {
  local dockerfile=$1
  local context=$2
  local tag=$3
  local outfile=$4

  echo "Building image ${tag} (Dockerfile: ${dockerfile}, context: ${context})"
  docker buildx build \
    -f "${dockerfile}" "${context}" \
    -t "${tag}" \
    --cache-from=type=local,src=${CACHE_DIR} \
    --cache-to=type=local,dest=${CACHE_DIR_NEW},mode=max \
    -o type=docker,dest="${outfile}",compression-level=0
}

# Run builds in parallel
build_image "Dockerfile" "." "${CONTROLLER_IMG_TAG}" "${DOCKER_IMAGES_PATH}/${CONTROLLER_IMG}-${GITHUB_SHA}" &
build_image "localmodel.Dockerfile" "." "${LOCALMODEL_CONTROLLER_IMG_TAG}" "${DOCKER_IMAGES_PATH}/${LOCALMODEL_CONTROLLER_IMG}-${GITHUB_SHA}" &
build_image "localmodel-agent.Dockerfile" "." "${LOCALMODEL_AGENT_IMG_TAG}" "${DOCKER_IMAGES_PATH}/${LOCALMODEL_AGENT_IMG}-${GITHUB_SHA}" &
build_image "agent.Dockerfile" "." "${AGENT_IMG_TAG}" "${DOCKER_IMAGES_PATH}/${AGENT_IMG}-${GITHUB_SHA}" &
build_image "router.Dockerfile" "." "${ROUTER_IMG_TAG}" "${DOCKER_IMAGES_PATH}/${ROUTER_IMG}-${GITHUB_SHA}" &

# Storage initializer (different context: python)
pushd python >/dev/null
  build_image "storage-initializer.Dockerfile" "." "${STORAGE_INIT_IMG_TAG}" "${DOCKER_IMAGES_PATH}/${STORAGE_INIT_IMG}-${GITHUB_SHA}" &
popd

wait  # wait for all parallel builds to finish

echo "Disk usage after building images:"
df -hT

# Rotate cache (preserve across builds)
rm -rf "${CACHE_DIR}"
mv "${CACHE_DIR_NEW}" "${CACHE_DIR}"

echo "✅ Done building all images"
