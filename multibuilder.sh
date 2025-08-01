#!/bin/bash

if [ -z "${MSC_VERSION}" ]; then
  echo "Error: MSC_VERSION is not set."
  exit 1
fi

GOOS="linux" GOARCH="amd64" go build -ldflags "-w -s -X 'github.com/monktype/msc/cmd.version=${MSC_VERSION}'" -o msc-${MSC_VERSION}_linux-amd64 .
GOOS="windows" GOARCH="amd64" go build -ldflags "-w -s -X 'github.com/monktype/msc/cmd.version=${MSC_VERSION}'" -o msc-${MSC_VERSION}_windows-amd64.exe .
GOOS="darwin" GOARCH="amd64" go build -ldflags "-w -s -X 'github.com/monktype/msc/cmd.version=${MSC_VERSION}'" -o msc-${MSC_VERSION}_macos-amd64 .
GOOS="darwin" GOARCH="arm64" go build -ldflags "-w -s -X 'github.com/monktype/msc/cmd.version=${MSC_VERSION}'" -o msc-${MSC_VERSION}_macos-arm64 .
