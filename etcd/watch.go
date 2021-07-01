package etcd

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

var HostMap map[string]map[string]int64
var HostSync map[string]sync.Once

func init() {
	HostMap = make(map[string]map[string]int64)
	HostSync = make(map[string]sync.Once)
}

func WatchList(keyPrefix string, ttl int64) {
	if once, ok := HostSync[keyPrefix]; !ok {
		HostSync[keyPrefix] = sync.Once{}
		once.Do(func() {
			rch := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrevKV())
			for n := range rch {
				for _, ev := range n.Events {
					switch ev.Type {
					case mvccpb.PUT:
						if _, ok := HostMap[keyPrefix]; !ok {
							HostMap[keyPrefix] = make(map[string]int64)
						}
						HostMap[keyPrefix][string(ev.Kv.Value)] = time.Now().Unix() + ttl + 1
					case mvccpb.DELETE:
						delete(HostMap,keyPrefix)
					}

				}
			}
		})
	}
}
