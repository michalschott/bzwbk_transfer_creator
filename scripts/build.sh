#!/usr/bin/env bash
#
# This script builds the application from source.

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd "$DIR"

# Get the git commit
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# Delete the old dir
echo "==> Removing old directory..."
rm -f bin/*
rm -rf pkg/*
mkdir -p bin/

# Allow LD_FLAGS to be appended during development compilations
LD_FLAGS="-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} $LD_FLAGS"

# In release mode we don't want debug information in the binary
if [[ -n "${RELEASE}" ]]; then
    LD_FLAGS="-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/michalschott/bzwbk_transfer_creator/bzwbk_transfer_creator.VersionPrerelease= -s -w"
fi

# Build!
echo "==> Building..."
go build \
    -o bin/bzwbk_transfer_creator \
    -ldflags "${LD_FLAGS}"

# Done!
echo
echo "==> Results:"
ls -hl bin/
