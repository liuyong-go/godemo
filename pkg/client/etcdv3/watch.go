package etcdv3

import (
	"context"
	"fmt"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/liuyong-go/godemo/pkg/util/ylog"
)

type Watch struct {
	revision     int64
	cancel       context.CancelFunc
	eventChan    chan *clientv3.Event
	lock         *sync.RWMutex
	incipientKVs []*mvccpb.KeyValue
}

//C 返回监听事件
func (w *Watch) C() chan *clientv3.Event {
	return w.eventChan
}

//IncipientKeyValues 返回incipient key and values
func (w *Watch) IncipientKeyValues() []*mvccpb.KeyValue {
	return w.incipientKVs
}
func (w *Watch) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}
func (client *Client) WatchPrefix(ctx context.Context, prefix string, fn func(ev *clientv3.Event) error) (*Watch, error) {
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var w = &Watch{
		revision:     resp.Header.Revision,
		eventChan:    make(chan *clientv3.Event, 100),
		incipientKVs: resp.Kvs,
	}
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		w.cancel = cancel
		rch := client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
		for {
			for n := range rch {
				if n.CompactRevision > w.revision {
					w.revision = n.CompactRevision
				}
				if n.Header.GetRevision() > w.revision {
					w.revision = n.Header.GetRevision()
				}
				if err := n.Err(); err != nil {
					ylog.SugarLogger.Error("etcd watch prefix", prefix, "error", err)
					continue
				}
				for _, ev := range n.Events {
					select {
					case w.eventChan <- ev:
						fmt.Println("biangeng", ev)
						_ = fn(ev)
					default:
						ylog.SugarLogger.Error("watch etcd with prefix block event chan, drop event message")
					}
				}
				ctx, cancel := context.WithCancel(context.Background())
				w.cancel = cancel
				if w.revision > 0 {
					rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
				} else {
					rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
				}

			}
		}
	}()
	return w, nil
}
