#! /bin/bash
set -x

go build
zip socksserver_linux64.zip socksserver

GOOS=darwin GOARCH=amd64 go build
zip socksserver_mac.zip socksserver

GOOS=windows GOARCH=amd64 go build
zip socksserver_windows64.zip socksserver.exe

GOOS=linux GOARCH=mipsle go build
zip socksserver_mipsle.zip socksserver

GOOS=linux GOARCH=arm go build
zip socksserver_arm.zip socksserver

