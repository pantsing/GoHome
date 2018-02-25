#!/usr/bin/env bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ "$GOPATH" == "" ]; then
    export GOPATH=$DIR
else
    export GOPATH=$DIR:$GOPATH
fi

GOFLAGS='-ldflags="-s -w"'
arch=$(go env GOARCH)
if [ "$2" == "arm" ]; then
    export GOARM="5"
    arch='arm'
fi
version=$(date +'%Y%m%d-%H%M%S')
goversion=$(go version | awk '{print $3}')
os='linux'
if [ "$1" == "darwin" ]; then
    os='darwin'
fi

GOOS=$os GOARCH=$arch GOFLAGS="$GOFLAGS" GOARM="${GOARM}" go build

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