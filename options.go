package service

import (
	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/client"
)

type clientKey struct{}

// Client to call broker service
func Client(c client.Client) broker.Option {
	return broker.SetOption(clientKey{}, c)
}

type serviceKey struct{}

// Service to call broker service
func Service(name string) broker.Option {
	return broker.SetOption(serviceKey{}, name)
}
