#!/bin/sh -eux

dep ensure
go build -o=sdchat.out "$@" ./cmd
