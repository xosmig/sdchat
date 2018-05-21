#!/bin/sh -eu

dep ensure
go build -o sdchat.out ./cmd
