#!/bin/sh

set -e

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

echo ${tag}
