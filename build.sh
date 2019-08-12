#!/usr/bin/env bash

go build -ldflags "-s -w" -o "bin/DbServiceWebApi" "src/main.go"
