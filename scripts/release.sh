#!/bin/bash
set -ex

PACKAGE_NAME=alox.sh

if [[ "$1" =~ ^[0-9]+(\.[0-9]+){2}$ ]]; then
    TAG_NAME="v$1"
fi

if [[ "$TAG_NAME" == "" ]]; then
    echo "Invalid version '$1'"
    echo "Expected version number in format 'x.y.z', where"
    echo "  x = major version"
    echo "  y = minor version"
    echo "  z = fix version"
    exit 1
fi

# ORIGINAL_BRANCH=$(git symbolic-ref --short HEAD)
# git checkout master
# git pull origin master

go mod tidy

go test ./...

# Reverting a remote git-anything is pain in
# the ass, better to be sure about the version
echo ""
echo "You are about to release a package version $TAG_NAME"
echo "Are you sure?"
read SURE

if [[ $SURE != "yes" ]]; then
    exit 0
fi

echo "really, Really sure?"
read SURE

if [[ $SURE != "yes" ]]; then
    exit 0
fi

git tag $TAG_NAME
git push origin $TAG_NAME

# git checkout $ORIGINAL_BRANCH

echo "Tagged latest master as version $TAG_NAME"

GOPROXY=proxy.golang.org go list -m $PACKAGE_NAME@$TAG_NAME

echo "Triggered indexing of package $PACKAGE_NAME@$TAG_NAME at proxy.golang.org"
