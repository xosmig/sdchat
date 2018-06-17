[![Build Status](https://travis-ci.org/xosmig/sdchat.svg?branch=master)](https://travis-ci.org/xosmig/sdchat)
[![Go Report Card](https://goreportcard.com/badge/github.com/xosmig/sdchat)](https://goreportcard.com/report/github.com/xosmig/sdchat)

# sdchat

Software design home assignment. Simple grpc one-on-one chat.

## Building and running

* install go: https://golang.org

* install dep: https://github.com/golang/dep

* run `./build.sh`

* run `./sdchat.out -port 8080 usernameServer` in one terminal

* run `./sdchat.out -serverip localhost -port 8080 usernameClient` in another terminal

* you can also run `./sdchat -help` for more detailed usage instructions

## Working with sources

* You can run `./test.sh` to run unit tests

* You can run `./generate.sh` to regenerate the generated sources (such as mock objects)
