#!/bin/bash

 

set -e

USER="{{ .User }}"
PASSWORD="{{ .Password }}"

SCHEME="{{ .Scheme }}"
REPO_HOST="{{ .Host }}"
REPO_PATH="{{ .Path }}"
REPO_NAME="{{ .Name }}"

[ "$1" = "--force" ] && FORCE=1

ARGS="repo add ${REPO_NAME} ${SCHEME}://${REPO_HOST}/${REPO_PATH}"

if ! which helm >/dev/null 2>&1; then
    echo "helm is required to setup the repository."
    exit 1
fi

if [[ -n "${USER}" ]]; then
    ARGS="${ARGS} --username ${USER}"
fi

if [[ -n "${PASSWORD}" ]]; then
    ARGS="${ARGS} --password ${PASSWORD}"
fi

if [[ -n "${FORCE}" ]]; then
    ARGS="${ARGS} --force-update"
fi

helm ${ARGS}
