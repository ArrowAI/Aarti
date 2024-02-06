#!/bin/bash

 

set -e

USER="{{ .User }}"
PASSWORD="{{ .Password }}"

SCHEME="{{ .Scheme }}"
REPO_HOST="{{ .Host }}"
REPO_PATH="{{ .Path }}"
REPO_NAME="{{ .Name }}"

DIST="{{ .Dist }}"
COMPONENT="{{ .Component }}"

[ -n "${DEBUG}" ] && set -x

if [ "$(id -u)" -ne 0 ]; then
    echo "Please run as root or sudo"
    exit 1
fi

if [ -f "/etc/apt/sources.list.d/${REPO_NAME}.list" ] && [ "$1" != "--force" ]; then
    echo "Repository already configured."
    echo "Use --force to overwrite."
    exit 1
fi

if ! which curl >/dev/null 2>&1; then
    echo "curl is required to setup the repository."
    exit 1
fi

REPO="${SCHEME}://${REPO_HOST}${REPO_PATH}"

if [ -n "${USER}" ]; then
    REPO_AUTH="${SCHEME}://${USER}:${PASSWORD}@${REPO_HOST}${REPO_PATH}"
    echo "machine ${REPO} login $USER password $PASSWORD" > "/etc/apt/auth.conf.d/${REPO_NAME}.conf"
else
    REPO_AUTH="${SCHEME}://${REPO_HOST}${REPO_PATH}"
fi

curl -s "${REPO_AUTH}/repository.key" -o "/etc/apt/trusted.gpg.d/${REPO_NAME}.asc"
echo "deb ${REPO} ${DIST} ${COMPONENT}" > "/etc/apt/sources.list.d/${REPO_NAME}.list"

echo "deb repository setup complete."
echo "You can now run 'apt update' to update the package list."

