#!/usr/bin/env bash

set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ "$GOPATH" == "" ]; then
    export GOPATH=$DIR
else
    export GOPATH=$DIR:$GOPATH
fi

GOFLAGS='-ldflags="-s -w"'
arch=$(go env GOARCH)
version=$(date +'%Y%m%d-%H%M%S')
goversion=$(go version | awk '{print $3}')
os='linux'
if [ "$1" == "darwin" ]; then
    os='darwin'
fi

GOOS=$os GOARCH=$arch GOFLAGS="$GOFLAGS" go build

if [ "$os" != "linux" ]; then
    exit
fi
TARGET=$DIR/dist/gohome
mkdir -p $TARGET
cp cert.pem key.pem gohome gohome.toml README.md systemd/* $TARGET
pushd $DIR/dist
tar -czvf $DIR/gh.tar.gz gohome
popd
rm -rf $DIR/dist/