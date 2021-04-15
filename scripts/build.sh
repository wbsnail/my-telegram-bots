#!/bin/sh

set -e
set -x

nodeRole=$1
commit=$2
tag=$3
preview=$4

if [ -z $preview ]; then
  namespace="default"
  env="production"
  logLevel="debug"
else
  namespace="preview"
  env="preview"
  logLevel="info"
fi

# working directory
WD=$(dirname "$(cd "$(dirname "$0")" && pwd)")
cd "${WD}" || echo "Cannot get into working directory: ${WD}"

docker build -t "my-telegram-bots:${tag}" .

sed "s/__namespace__/${namespace}/g;
     s/__env__/${env}/g;
     s/__log_level__/${logLevel}/g;
     s/__tag__/${tag}/g;
     s/__node_role__/${nodeRole}/g;
     s/commit_placeholder/${commit}/g;" \
  ./deploy/deployment.tmpl.yaml >./deploy/deployment.yaml

kubectl apply -f ./deploy/deployment.yaml
