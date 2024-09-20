package service

import (
	"context"

	pbmicro "go.unistack.org/micro-broker-service/v3/micro"
	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
)

type serviceSub struct {
	topic   string
	group   string
	handler broker.Handler
	stream  pbmicro.BrokerService_SubscribeClient
	closed  chan bool
	options broker.SubscribeOptions
}

type serviceEvent struct {
	ctx     context.Context
	topic   string
	err     error
	message *broker.Message
}

func (s *serviceEvent) Context() context.Context {
	return s.ctx
}

func (s *serviceEvent) Topic() string {
	return s.topic
}

func (s *serviceEvent) Message() *broker.Message {
	return s.message
}

func (s *serviceEvent) Ack() error {
	return nil
}

func (s *serviceEvent) Error() error {
	return s.err
}

func (s *serviceEvent) SetError(err error) {
	s.err = err
}

func (s *serviceSub) isClosed() bool {
	select {
	case <-s.closed:
		return true
	default:
		return false
	}
}

func (s *serviceSub) run(ctx context.Context) error {
	exit := make(chan bool)
	go func() {
		select {
		case <-exit:
		case <-s.closed:
		}

		// close the stream
		s.stream.Close()
	}()

	for {
		// TODO: do not fail silently
		msg, err := s.stream.Recv()
		ctx := metadata.NewIncomingContext(context.Background(), msg.Header)
		if err != nil {
			if logger.V(logger.TraceLevel) {
				logger.Tracef(ctx, "streaming error for subcription to topic %s: %v", s.Topic(), err)
			}

			// close the exit channel
			close(exit)

			// don't return an error if we unsubscribed
			if s.isClosed() {
				return nil
			}

			// return stream error
			return err
		}

		p := &serviceEvent{
			ctx:   ctx,
			topic: s.topic,
			message: &broker.Message{
				Header: msg.Header,
				Body:   msg.Body,
			},
		}
		p.err = s.handler(p)
	}
}

func (s *serviceSub) Options() broker.SubscribeOptions {
	return s.options
}

func (s *serviceSub) Topic() string {
	return s.topic
}

func (s *serviceSub) Unsubscribe(ctx context.Context) error {
	select {
	case <-s.closed:
		return nil
	default:
		close(s.closed)
	}
	return nil
}
