// Package service provides the broker service client
package service

import (
	"context"
	"fmt"
	"time"

	pbmicro "github.com/unistack-org/micro-broker-service/v3/micro"
	pb "github.com/unistack-org/micro-broker-service/v3/proto"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/logger"
)

type serviceBroker struct {
	addrs   []string
	service string
	client  pbmicro.BrokerClient
	init    bool
	opts    broker.Options
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

	b.client = pbmicro.NewBrokerClient(b.service, cli)

	return nil
}

func (b *serviceBroker) Options() broker.Options {
	return b.opts
}

func (b *serviceBroker) Publish(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	if logger.V(logger.TraceLevel) {
		logger.Tracef(ctx, "Publishing to topic %s broker %v", topic, b.addrs)
	}
	_, err := b.client.Publish(ctx, &pb.PublishRequest{
		Topic: topic,
		Message: &pb.Message{
			Header: msg.Header,
			Body:   msg.Body,
		},
	}, client.WithAddress(b.addrs...))
	return err
}

func (b *serviceBroker) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	options := broker.NewSubscribeOptions(opts...)
	if logger.V(logger.TraceLevel) {
		logger.Tracef(ctx, "Subscribing to topic %s group %s broker %v", topic, options.Group, b.addrs)
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
	}

	go func() {
		for {
			select {
			case <-sub.closed:
				if logger.V(logger.TraceLevel) {
					logger.Tracef(ctx, "Unsubscribed from topic %s", topic)
				}
				return
			default:
				if logger.V(logger.TraceLevel) {
					// run the subscriber
					logger.Tracef(ctx, "Streaming from broker %v to topic [%s] group [%s]", b.addrs, topic, options.Group)
				}
				if err := sub.run(ctx); err != nil {
					if logger.V(logger.TraceLevel) {
						logger.Tracef(ctx, "Resubscribing to topic %s broker %v", topic, b.addrs)
					}
					stream, err := b.client.Subscribe(ctx, &pb.SubscribeRequest{
						Topic: topic,
						Group: options.Group,
					}, client.WithAddress(b.addrs...), client.WithRequestTimeout(time.Hour))
					if err != nil {
						if logger.V(logger.TraceLevel) {
							logger.Tracef(ctx, "Failed to resubscribe to topic %s: %v", topic, err)
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
