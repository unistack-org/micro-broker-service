// Package service provides the broker service client
package service // import "go.unistack.org/micro-broker-service/v3"

import (
	"context"
	"fmt"
	"time"

	pbmicro "go.unistack.org/micro-broker-service/v3/micro"
	pb "go.unistack.org/micro-broker-service/v3/proto"
	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
)

var _ broker.Broker = &serviceBroker{}

type serviceBroker struct {
	client  pbmicro.BrokerServiceClient
	service string
	opts    broker.Options
	addrs   []string
	init    bool
}

func (b *serviceBroker) Address() string {
	return b.addrs[0]
}

func (b *serviceBroker) Connect(ctx context.Context) error {
	return nil
}

func (b *serviceBroker) Disconnect(ctx context.Context) error {
	return nil
}

func (b *serviceBroker) Init(opts ...broker.Option) error {
	if len(opts) == 0 && b.init {
		return nil
	}

	for _, o := range opts {
		o(&b.opts)
	}

	var cli client.Client
	if b.opts.Context != nil {
		if v, ok := b.opts.Context.Value(clientKey{}).(client.Client); ok && v != nil {
			cli = v
		}
		if v, ok := b.opts.Context.Value(serviceKey{}).(string); ok && v != "" {
			b.service = v
		}
	}

	if err := b.opts.Register.Init(); err != nil {
		return err
	}
	if err := b.opts.Tracer.Init(); err != nil {
		return err
	}
	if err := b.opts.Logger.Init(); err != nil {
		return err
	}
	if err := b.opts.Meter.Init(); err != nil {
		return err
	}

	if b.service == "" {
		return fmt.Errorf("missing Service option")
	}

	if cli == nil {
		return fmt.Errorf("missing Client option")
	}

	b.client = pbmicro.NewBrokerServiceClient(b.service, cli)

	b.init = true
	return nil
}

func (b *serviceBroker) Options() broker.Options {
	return b.opts
}

func (b *serviceBroker) BatchPublish(ctx context.Context, msgs []*broker.Message, opts ...broker.PublishOption) error {
	return b.publish(ctx, msgs, opts...)
}

func (b *serviceBroker) Publish(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	msg.Header.Set(metadata.HeaderTopic, topic)
	return b.publish(ctx, []*broker.Message{msg}, opts...)
}

func (b *serviceBroker) publish(ctx context.Context, msgs []*broker.Message, opts ...broker.PublishOption) error {
	for _, msg := range msgs {
		topic, _ := msg.Header.Get(metadata.HeaderTopic)
		if b.opts.Logger.V(logger.TraceLevel) {
			b.opts.Logger.Trace(ctx, fmt.Sprintf("Publishing to topic %s broker %v", topic, b.addrs))
		}
		_, err := b.client.Publish(ctx, &pb.PublishRequest{
			Topic: topic,
			Message: &pb.Message{
				Header: msg.Header,
				Body:   msg.Body,
			},
		}, client.WithAddress(b.addrs...))
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *serviceBroker) BatchSubscribe(ctx context.Context, topic string, handler broker.BatchHandler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return nil, nil
}

func (b *serviceBroker) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	options := broker.NewSubscribeOptions(opts...)
	if b.opts.Logger.V(logger.TraceLevel) {
		b.opts.Logger.Trace(ctx, fmt.Sprintf("Subscribing to topic %s group %s broker %v", topic, options.Group, b.addrs))
	}
	stream, err := b.client.Subscribe(ctx, &pb.SubscribeRequest{
		Topic: topic,
		Group: options.Group,
	}, client.WithAddress(b.addrs...), client.WithRequestTimeout(time.Hour))
	if err != nil {
		return nil, err
	}

	sub := &serviceSub{
		topic:   topic,
		group:   options.Group,
		handler: handler,
		stream:  stream,
		closed:  make(chan bool),
		options: options,
		opts:    b.opts,
	}

	go func() {
		for {
			select {
			case <-sub.closed:
				if b.opts.Logger.V(logger.TraceLevel) {
					b.opts.Logger.Trace(ctx, "Unsubscribed from topic %s", topic)
				}
				return
			default:
				if b.opts.Logger.V(logger.TraceLevel) {
					// run the subscriber
					b.opts.Logger.Trace(ctx, fmt.Sprintf("Streaming from broker %v to topic [%s] group [%s]", b.addrs, topic, options.Group))
				}
				if err := sub.run(ctx); err != nil {
					if b.opts.Logger.V(logger.TraceLevel) {
						b.opts.Logger.Trace(ctx, fmt.Sprintf("Resubscribing to topic %s broker %v", topic, b.addrs))
					}
					stream, err := b.client.Subscribe(ctx, &pb.SubscribeRequest{
						Topic: topic,
						Group: options.Group,
					}, client.WithAddress(b.addrs...), client.WithRequestTimeout(time.Hour))
					if err != nil {
						if b.opts.Logger.V(logger.TraceLevel) {
							b.opts.Logger.Trace(ctx, fmt.Sprintf("Failed to resubscribe to topic %s: %v", topic, err))
						}
						time.Sleep(time.Second)
						continue
					}
					// new stream
					sub.stream = stream
				}
			}
		}
	}()

	return sub, nil
}

func (b *serviceBroker) String() string {
	return "service"
}

func (b *serviceBroker) Name() string {
	return b.opts.Name
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.NewOptions(opts...)

	addrs := options.Addrs
	if len(addrs) == 0 {
		addrs = []string{"127.0.0.1:8001"}
	}

	return &serviceBroker{
		addrs: addrs,
		opts:  options,
	}
}
