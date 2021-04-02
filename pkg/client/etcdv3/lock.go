package etcdv3

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3/concurrency"
)

type Mutex struct {
	s *concurrency.Session
	m *concurrency.Mutex
}

func (client *Client) NewMutex(key string, opts ...concurrency.SessionOption) (mutex *Mutex, err error) {
	mutex = &Mutex{}
	mutex.s, err = concurrency.NewSession(client.Client, opts...)
	if err != nil {
		return
	}
	mutex.m = concurrency.NewMutex(mutex.s, key)
	return
}

// Lock ...
func (mutex *Mutex) Lock(timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return mutex.m.Lock(ctx)
}
func (mutex *Mutex) Unlock() (err error) {
	err = mutex.m.Unlock(context.TODO())
	if err != nil {
		return
	}
	return mutex.s.Close()
}
