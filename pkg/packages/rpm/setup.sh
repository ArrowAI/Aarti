#!/bin/bash

 

set -e

USER="{{ .User }}"
PASSWORD="{{ .Password }}"

SCHEME="{{ .Scheme }}"
REPO_HOST="{{ .Host }}"
REPO_PATH="{{ .Path }}"

REPO_NAME="{{ .Name }}"

already_exists() {
    echo "Repository already configured."
    echo "Use --force to overwrite."
}

[ -n "${DEBUG}" ] && set -x

if  [ -z "${REPO_PATH}" ]; then
    REPO_PATH="/"
fi

if [ -n "${USER}" ]; then
    REPO_URL="${SCHEME}://${USER}:${PASSWORD}@${REPO_HOST}${REPO_PATH}"
else
    REPO_URL="${SCHEME}://${REPO_HOST}${REPO_PATH}"
fi

if [ "$UID" -ne 0 ]; then
    echo "Please run as root or sudo"
    exit 1
fi

if [ -f "/etc/yum.repos.d/${REPO_NAME}.repo" ] && [ "$1" != "--force" ]; then
    echo "Repository already configured."
    echo "Use --force to overwrite."
    exit 1
fi

if ! command -v curl > /dev/null; then
    echo "curl is required to setup the repository."
    exit 1
fi

curl -s "${REPO_URL}.repo" -o "/etc/yum.repos.d/${REPO_NAME}.repo"

echo "yum setup complete."
echo "You can now run 'yum update' to update the package list."
