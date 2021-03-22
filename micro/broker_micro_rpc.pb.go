// Code generated by protoc-gen-micro
// source: broker.proto
package service

import (
	context "context"
	proto "github.com/unistack-org/micro-broker-service/v3/proto"
	api "github.com/unistack-org/micro/v3/api"
	client "github.com/unistack-org/micro/v3/client"
	server "github.com/unistack-org/micro/v3/server"
)

type brokerClient struct {
	c    client.Client
	name string
}

func NewBrokerClient(name string, c client.Client) BrokerClient {
	return &brokerClient{c: c, name: name}
}

func (c *brokerClient) Publish(ctx context.Context, req *proto.PublishRequest, opts ...client.CallOption) (*proto.Empty, error) {
	rsp := &proto.Empty{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Broker.Publish", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *brokerClient) Subscribe(ctx context.Context, req *proto.SubscribeRequest, opts ...client.CallOption) (Broker_SubscribeClient, error) {
	stream, err := c.c.Stream(ctx, c.c.NewRequest(c.name, "Broker.Subscribe", &proto.SubscribeRequest{}), opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(req); err != nil {
		return nil, err
	}
	return &brokerClientSubscribe{stream}, nil
}

type brokerClientSubscribe struct {
	stream client.Stream
}

func (s *brokerClientSubscribe) Close() error {
	return s.stream.Close()
}

func (s *brokerClientSubscribe) Context() context.Context {
	return s.stream.Context()
}

func (s *brokerClientSubscribe) SendMsg(msg interface{}) error {
	return s.stream.Send(msg)
}

func (s *brokerClientSubscribe) RecvMsg(msg interface{}) error {
	return s.stream.Recv(msg)
}

func (s *brokerClientSubscribe) Recv() (*proto.Message, error) {
	msg := &proto.Message{}
	if err := s.stream.Recv(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

type brokerServer struct {
	BrokerServer
}

func (h *brokerServer) Publish(ctx context.Context, req *proto.PublishRequest, rsp *proto.Empty) error {
	return h.BrokerServer.Publish(ctx, req, rsp)
}

func (h *brokerServer) Subscribe(ctx context.Context, stream server.Stream) error {
	msg := &proto.SubscribeRequest{}
	if err := stream.Recv(msg); err != nil {
		return err
	}
	return h.BrokerServer.Subscribe(ctx, msg, &brokerSubscribeStream{stream})
}

type brokerSubscribeStream struct {
	stream server.Stream
}

func (s *brokerSubscribeStream) Close() error {
	return s.stream.Close()
}

func (s *brokerSubscribeStream) Context() context.Context {
	return s.stream.Context()
}

func (s *brokerSubscribeStream) SendMsg(msg interface{}) error {
	return s.stream.Send(msg)
}

func (s *brokerSubscribeStream) RecvMsg(msg interface{}) error {
	return s.stream.Recv(msg)
}

func (s *brokerSubscribeStream) Send(msg *proto.Message) error {
	return s.stream.Send(msg)
}

func RegisterBrokerServer(s server.Server, sh BrokerServer, opts ...server.HandlerOption) error {
	type broker interface {
		Publish(ctx context.Context, req *proto.PublishRequest, rsp *proto.Empty) error
		Subscribe(ctx context.Context, stream server.Stream) error
	}
	type Broker struct {
		broker
	}
	h := &brokerServer{sh}
	for _, endpoint := range NewBrokerEndpoints() {
		opts = append(opts, api.WithEndpoint(endpoint))
	}
	return s.Handle(s.NewHandler(&Broker{h}, opts...))
}