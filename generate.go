package service

//go:generate protoc -I./proto -I. --go-grpc_out=paths=source_relative:./proto --go_out=paths=source_relative:./proto --micro_out=components=micro|grpc,paths=source_relative:./proto proto/broker.proto
