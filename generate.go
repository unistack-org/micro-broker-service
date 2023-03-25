package service

//go:generate protoc -I./proto -I. --go_out=paths=source_relative:./proto --go-micro_out=components=micro|grpc,standalone=true,paths=source_relative:./micro proto/broker.proto
