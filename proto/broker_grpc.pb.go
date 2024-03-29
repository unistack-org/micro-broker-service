// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package servicepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BrokerServiceClient is the client API for BrokerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrokerServiceClient interface {
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	BatchPublish(ctx context.Context, in *BatchPublishRequest, opts ...grpc.CallOption) (*BatchPublishResponse, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (BrokerService_SubscribeClient, error)
	BatchSubscribe(ctx context.Context, in *BatchSubscribeRequest, opts ...grpc.CallOption) (BrokerService_BatchSubscribeClient, error)
}

type brokerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBrokerServiceClient(cc grpc.ClientConnInterface) BrokerServiceClient {
	return &brokerServiceClient{cc}
}

func (c *brokerServiceClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, "/servicepb.BrokerService/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerServiceClient) BatchPublish(ctx context.Context, in *BatchPublishRequest, opts ...grpc.CallOption) (*BatchPublishResponse, error) {
	out := new(BatchPublishResponse)
	err := c.cc.Invoke(ctx, "/servicepb.BrokerService/BatchPublish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (BrokerService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &BrokerService_ServiceDesc.Streams[0], "/servicepb.BrokerService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &brokerServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BrokerService_SubscribeClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type brokerServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *brokerServiceSubscribeClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *brokerServiceClient) BatchSubscribe(ctx context.Context, in *BatchSubscribeRequest, opts ...grpc.CallOption) (BrokerService_BatchSubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &BrokerService_ServiceDesc.Streams[1], "/servicepb.BrokerService/BatchSubscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &brokerServiceBatchSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BrokerService_BatchSubscribeClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type brokerServiceBatchSubscribeClient struct {
	grpc.ClientStream
}

func (x *brokerServiceBatchSubscribeClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BrokerServiceServer is the server API for BrokerService service.
// All implementations must embed UnimplementedBrokerServiceServer
// for forward compatibility
type BrokerServiceServer interface {
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	BatchPublish(context.Context, *BatchPublishRequest) (*BatchPublishResponse, error)
	Subscribe(*SubscribeRequest, BrokerService_SubscribeServer) error
	BatchSubscribe(*BatchSubscribeRequest, BrokerService_BatchSubscribeServer) error
	mustEmbedUnimplementedBrokerServiceServer()
}

// UnimplementedBrokerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBrokerServiceServer struct {
}

func (UnimplementedBrokerServiceServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedBrokerServiceServer) BatchPublish(context.Context, *BatchPublishRequest) (*BatchPublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchPublish not implemented")
}
func (UnimplementedBrokerServiceServer) Subscribe(*SubscribeRequest, BrokerService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedBrokerServiceServer) BatchSubscribe(*BatchSubscribeRequest, BrokerService_BatchSubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method BatchSubscribe not implemented")
}
func (UnimplementedBrokerServiceServer) mustEmbedUnimplementedBrokerServiceServer() {}

// UnsafeBrokerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrokerServiceServer will
// result in compilation errors.
type UnsafeBrokerServiceServer interface {
	mustEmbedUnimplementedBrokerServiceServer()
}

func RegisterBrokerServiceServer(s grpc.ServiceRegistrar, srv BrokerServiceServer) {
	s.RegisterService(&BrokerService_ServiceDesc, srv)
}

func _BrokerService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/servicepb.BrokerService/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServiceServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BrokerService_BatchPublish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchPublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServiceServer).BatchPublish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/servicepb.BrokerService/BatchPublish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServiceServer).BatchPublish(ctx, req.(*BatchPublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BrokerService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BrokerServiceServer).Subscribe(m, &brokerServiceSubscribeServer{stream})
}

type BrokerService_SubscribeServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type brokerServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *brokerServiceSubscribeServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func _BrokerService_BatchSubscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BatchSubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BrokerServiceServer).BatchSubscribe(m, &brokerServiceBatchSubscribeServer{stream})
}

type BrokerService_BatchSubscribeServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type brokerServiceBatchSubscribeServer struct {
	grpc.ServerStream
}

func (x *brokerServiceBatchSubscribeServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

// BrokerService_ServiceDesc is the grpc.ServiceDesc for BrokerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BrokerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "servicepb.BrokerService",
	HandlerType: (*BrokerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _BrokerService_Publish_Handler,
		},
		{
			MethodName: "BatchPublish",
			Handler:    _BrokerService_BatchPublish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _BrokerService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BatchSubscribe",
			Handler:       _BrokerService_BatchSubscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "broker.proto",
}
