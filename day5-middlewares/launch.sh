#!/bin/bash
go build -o gee main.go
dlv debug --headless --listen=:2345 --api-version=2 --log=true ./gee