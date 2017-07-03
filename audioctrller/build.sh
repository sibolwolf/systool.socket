#!/bin/sh

source /opt/toolchain-sunxi/environment-setup-arm-openwrt-linux

GOARCH=arm GOARM=7 CGO_ENABLED=1 go test -o record.test -c