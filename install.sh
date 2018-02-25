#!/usr/bin/env bash

tar -C /usr/local/ -xzf gh.tar.gz
cp /usr/local/gohome/*.service /usr/lib/systemd/system/
#systemctl enable ghserver.service
systemctl enable ghclient.service