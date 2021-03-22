package service

import (
	"context"

	pbmicro "github.com/unistack-org/micro-broker-service/v3/micro"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/logger"
)

type serviceSub struct {
	topic   string
	group   string
	handler broker.Handler
	stream  pbmicro.Broker_SubscribeClient
	closed  chan bool
	options broker.SubscribeOptions
}

type serviceEvent struct {
	topic   string
	err     error
	message *broker.Message
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
		if err != nil {
			if logger.V(logger.TraceLevel) {
				logger.Tracef(ctx, "Streaming error for subcription to topic %s: %v", s.Topic(), err)
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
