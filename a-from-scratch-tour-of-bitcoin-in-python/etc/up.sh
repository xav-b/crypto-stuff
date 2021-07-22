#! /usr/bin/env bash

# variables (
  _VENV="crypto"
  _ENV_FILE=".env"
# )

_log() {
  echo -e "$(date '+%H:%M:%S') [init] $@"
}

echo -e "\n======================================================================================\n"

_log "setting up project workspace for development\n"
echo -e "\tdir: $(PWD)"
echo -e "\tVCS: git::$(git rev-parse --abbrev-ref HEAD)"
echo

_log "runtime: activating virtualenv ${_VENV}"
pyenv activate "$_VENV"

_log "runtime: using $(python --version)"

export PYTHONPATH="$PWD:$PYTHONPATH"

_log "loading env vars from '${_ENV_FILE}'"
# https://stackoverflow.com/questions/19331497/set-environment-variables-from-file-of-key-value-pairs
set -o allexport
source "$_ENV_FILE"
set +o allexport

echo
_log "environment ready ✨"

echo -e "\n======================================================================================\n"
