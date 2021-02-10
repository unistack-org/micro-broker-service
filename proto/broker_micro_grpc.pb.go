// Code generated by protoc-gen-micro
// source: broker.proto
package service

import (
	"context"

	micro_client "github.com/unistack-org/micro/v3/client"
	micro_server "github.com/unistack-org/micro/v3/server"
)

var (
	_ micro_server.Option
	_ micro_client.Option
)

type brokerService struct {
	c    micro_client.Client
	name string
}

// Micro client stuff

// NewBrokerService create new service client
func NewBrokerService(name string, c micro_client.Client) BrokerService {
	return &brokerService{c: c, name: name}
}

func (c *brokerService) Publish(ctx context.Context, req *PublishRequest, opts ...micro_client.CallOption) (*Empty, error) {
	rsp := &Empty{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Broker.Publish", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *brokerService) Subscribe(ctx context.Context, req *SubscribeRequest, opts ...micro_client.CallOption) (Broker_SubscribeService, error) {
	stream, err := c.c.Stream(ctx, c.c.NewRequest(c.name, "Broker.Subscribe", &SubscribeRequest{}), opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(req); err != nil {
		return nil, err
	}
	return &brokerServiceSubscribe{stream}, nil
}

type brokerServiceSubscribe struct {
	stream micro_client.Stream
}

func (x *brokerServiceSubscribe) Close() error {
	return x.stream.Close()
}

func (x *brokerServiceSubscribe) Context() context.Context {
	return x.stream.Context()
}

func (x *brokerServiceSubscribe) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *brokerServiceSubscribe) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *brokerServiceSubscribe) Recv() (*Message, error) {
	m := &Message{}
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
} // Micro server stuff

type brokerHandler struct {
	BrokerHandler
}

func (h *brokerHandler) Publish(ctx context.Context, req *PublishRequest, rsp *Empty) error {
	return h.BrokerHandler.Publish(ctx, req, rsp)
}

func (h *brokerHandler) Subscribe(ctx context.Context, stream micro_server.Stream) error {
	m := &SubscribeRequest{}
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.BrokerHandler.Subscribe(ctx, m, &brokerSubscribeStream{stream})
}

type brokerSubscribeStream struct {
	stream micro_server.Stream
}

func (x *brokerSubscribeStream) Close() error {
	return x.stream.Close()
}

func (x *brokerSubscribeStream) Context() context.Context {
	return x.stream.Context()
}

func (x *brokerSubscribeStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *brokerSubscribeStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}
func (x *brokerSubscribeStream) Send(m *Message) error {
	return x.stream.Send(m)
}