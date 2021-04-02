package etcdv3

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/util/ylog"
	"google.golang.org/grpc"
)

var etcdCtl *Client

type Client struct {
	*clientv3.Client
	config *conf.EtcdV3
}

func NewClient() *Client {
	if etcdCtl != nil {
		return etcdCtl
	}
	config := conf.NewEtcdV3().GetConfig()
	conf := clientv3.Config{
		Endpoints:            config.Endpoints,
		DialTimeout:          config.ConnectTimeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		DialOptions: []grpc.DialOption{ //grpc 配置
			grpc.WithBlock(),
		},
		AutoSyncInterval: config.AutoSyncInterval,
	}
	if config.Endpoints == nil {
		ylog.SugarLogger.Panic("client etcd endpoints empty")
	}
	if !config.Secure {
		conf.DialOptions = append(conf.DialOptions, grpc.WithInsecure())
	}
	if config.BasicAuth {
		conf.Username = config.UserName
		conf.Password = config.Password
	}
	tlsEnabled := false
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}
	if config.CaCert != "" {
		certBytes, err := ioutil.ReadFile(config.CaCert)
		if err != nil {
			ylog.SugarLogger.Panic("parse CaCert failed", err)
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)

		if ok {
			tlsConfig.RootCAs = caCertPool
		}
		tlsEnabled = true
	}
	if config.CertFile != "" && config.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			ylog.SugarLogger.Panic("load CertFile or KeyFile failed", err)
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}
	if tlsEnabled {
		conf.TLS = tlsConfig
	}
	client, err := clientv3.New(conf)
	if err != nil {
		ylog.SugarLogger.Panic("client etcd start panic", err)
	}
	etcdCtl = &Client{
		Client: client,
		config: config,
	}

	return etcdCtl
}

// GetKeyValue queries etcd key, returns mvccpb.KeyValue
func (client *Client) GetKeyValue(ctx context.Context, key string) (kv *mvccpb.KeyValue, err error) {
	rp, err := client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(rp.Kvs) > 0 {
		return rp.Kvs[0], nil
	}
	return
}

// GetPrefix get prefix
func (client *Client) GetPrefix(ctx context.Context, prefix string) (map[string]string, error) {
	var vars = make(map[string]string)
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return vars, err
	}
	for _, kv := range resp.Kvs {
		vars[string(kv.Key)] = string(kv.Value)
	}
	return vars, nil
}

// DelPrefix 按前缀删除
func (client *Client) DelPrefix(ctx context.Context, prefix string) (deleted int64, err error) {
	resp, err := client.Delete(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return 0, err
	}
	return resp.Deleted, nil
}

//GetLeaseSession 创建租约会话
func (client *Client) GetLeaseSession(ctx context.Context, opts ...concurrency.SessionOption) (leaseSession *concurrency.Session, err error) {
	return concurrency.NewSession(client.Client, opts...)
}
func (client *Client) CloseLeaseSession(session *concurrency.Session) {
	err := session.Close()
	if err != nil {
		ylog.SugarLogger.Info("close session error", session)
	}
	return
}
