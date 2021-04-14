package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/liuyong-go/godemo/pkg/util/ylog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Option func(c *Client)
type Client struct {
	Adress      string
	Block       bool
	DialTimeout time.Duration
	ReadTimeout time.Duration
	KeepAlive   *keepalive.ClientParameters
	dialOptions []grpc.DialOption
}

func defaultClient() *Client {
	return &Client{
		dialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		DialTimeout: time.Second * 3,
		ReadTimeout: time.Second * 1,
		Block:       true,
	}
}
func WithAdress(address string) Option {
	return func(c *Client) {
		c.Adress = address
	}
}
func WithBlock(block bool) Option {
	return func(c *Client) {
		c.Block = block
	}
}
func WithDialTimeOut(t time.Duration) Option {
	return func(c *Client) {
		c.DialTimeout = t
	}
}
func ApplyOption(options ...Option) *Client {
	client := defaultClient()
	for _, option := range options {
		option(client)
	}
	return client
}
func NewGRPCClient(options ...Option) *grpc.ClientConn {
	config := ApplyOption(options...)
	var ctx = context.Background()
	var dialOptions = config.dialOptions
	if config.Block {
		if config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
			defer cancel()
		}
		dialOptions = append(dialOptions, grpc.WithBlock())
	}
	if config.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.KeepAlive))

	}
	cc, err := grpc.DialContext(ctx, config.Adress, dialOptions...)
	if err != nil {
		if ylog.Logger != nil {
			ylog.SugarLogger.Error("dial grpc server err", err)
		} else {
			fmt.Println("dial grpc server err", err)
		}

	}
	return cc

}
