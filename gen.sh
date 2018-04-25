#!/bin/sh
protoc -I proto proto/elog.proto --go_out=plugins=grpc:./
