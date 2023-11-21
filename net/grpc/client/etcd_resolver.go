package client

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	gresolver "google.golang.org/grpc/resolver"
	"strings"
)

type builder struct {
	c *clientv3.Client
}

func (b builder) Build(target gresolver.Target, cc gresolver.ClientConn, opts gresolver.BuildOptions) (gresolver.Resolver, error) {
	// Refer to https://github.com/grpc/grpc-go/blob/16d3df80f029f57cff5458f1d6da6aedbc23545d/clientconn.go#L1587-L1611
	endpoint := target.URL.Path
	s := strings.Split(endpoint, "/")
	resp, err := b.c.Get(context.Background(), s[1], clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	addrs := make([]gresolver.Address, 0, len(resp.Kvs))
	for _, v := range resp.Kvs {

		addr := gresolver.Address{
			Addr: string(v.Value),
		}
		addrs = append(addrs, addr)
	}
	if err := cc.UpdateState(gresolver.State{Addresses: addrs}); err != nil {
		return nil, err
	}

	return &resolver{}, nil
}

func (b builder) Scheme() string {
	return "etcd"
}

// NewBuilder creates a resolver builder.
func NewBuilder(client *clientv3.Client) (gresolver.Builder, error) {
	return builder{c: client}, nil
}

type resolver struct {
}

// ResolveNow is a no-op here.
// It's just a hint, resolver can ignore this if it's not necessary.
func (r *resolver) ResolveNow(gresolver.ResolveNowOptions) {}

func (r *resolver) Close() {

}
