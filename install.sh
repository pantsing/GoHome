#!/usr/bin/env bash

tar -C /usr/local/ gh.tar.gz
pushd /usr/local/gohome/
cp *.service /usr/lib/systemd/system/
popd
systemctl enable ghserver.service
systemctl enable ghclient.service