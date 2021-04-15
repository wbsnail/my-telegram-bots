#!/bin/sh

set -e
set -x

preview=$1

if [ -z $preview ]; then
  if [ -z ${CIRCLE_TAG} ]; then
    tag="unknown"
  else
    tag=${CIRCLE_TAG}
  fi
else
  tag="preview"
fi

echo 'Syncing files'

BUILD_DIR="${CI_DIR}/my-telegram-bots"

rsync \
  -avz \
  -e "ssh -o StrictHostKeyChecking=no -p ${SSH_PORT}" \
  --delete-after \
  --progress \
  --relative \
  --rsync-path="mkdir -p ${BUILD_DIR} && rsync" \
  build \
  deploy \
  docker_image \
  frontend/dist \
  scripts \
  templates \
  translations \
  vite/dist \
  "${SSH_USER}"@"${SSH_ADDRESS}":"${BUILD_DIR}"

echo 'Sync finished, building'

echo "Build command: ${BUILD_DIR}/scripts/build.sh ${CIRCLE_SHA1} ${tag} ${preview}"

ssh \
  -i ./credentials/id_rsa \
  -o UserKnownHostsFile=/dev/null \
  -o StrictHostKeyChecking=no \
  -p "${SSH_PORT}" \
  "${SSH_USER}"@"${SSH_ADDRESS}" "${BUILD_DIR}/scripts/build.sh ${NODE_ROLE} ${CIRCLE_SHA1} ${tag} ${preview}"
