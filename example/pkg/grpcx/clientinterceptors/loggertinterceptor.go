package clientinterceptors

import (
	"context"
	"github.com/luxun9527/zlog"
	"google.golang.org/grpc"
)

// 从go-zero中复制的
// LoggerCallOption is a call option that controls timeout.
type LoggerCallOption struct {
	grpc.EmptyCallOption
}

// LoggerInterceptor is an interceptor that controls timeout.
func LoggerInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		zlog.Debugf("[GRPC] client call method %s", method)
		return err
	}
}

// WithCallLogger returns a call option that controls method call timeout.
func WithCallLogger() grpc.CallOption {
	return LoggerCallOption{}
}
