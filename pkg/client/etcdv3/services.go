package etcdv3

import (
	"context"
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/liuyong-go/godemo/pkg/server"
)

//RegistService 注册服务并续约
func (client *Client) RegistService(s *server.ServiceInfo) error {
	ctx := context.Context(context.Background())
	sessionLease, err := client.GetLeaseSession(ctx, concurrency.WithTTL(10))
	if err != nil {
		fmt.Println("get lease session fail", err)
		return err
	}
	serviceKey := client.config.ServicePrefix + s.Scheme + "/" + s.Name
	_, err = client.Put(ctx, serviceKey, s.Address, clientv3.WithLease(sessionLease.Lease()))
	if err != nil {
		fmt.Println("register service fail", serviceKey, s.Address, sessionLease.Lease())
		return err
	}
	s.Session = sessionLease
	return nil
}
func (client *Client) UnregistService(s *server.ServiceInfo) error {
	fmt.Println("close session", s.Session.Lease())
	client.CloseLeaseSession(s.Session)
	return nil
}
