#!/usr/bin/env bash

go mod download
go test -v ./...
