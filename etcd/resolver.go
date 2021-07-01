package etcd

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"time"
)

type etcdResolver struct {
	rawAddr []string
	cc      resolver.ClientConn
	load    Load
	ttl     int64
}

func NewResolver(etcdAddr []string, load Load, ttl int64) resolver.Builder {
	return &etcdResolver{rawAddr: etcdAddr, load: load, ttl: ttl}
}

// Build 构建etcd client
func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   r.rawAddr,
			DialTimeout: 15 * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}
	r.cc = cc
	go r.watch("/" + target.Scheme + "/" + target.Endpoint)
	return r, nil
}

// watch 监听resolve列表变化
func (r *etcdResolver) watch(keyPrefix string) {
	var addrList []resolver.Address
	addr := r.load.getAddr(HostMap[keyPrefix])
	if addr != "" {
		addrList = append(addrList, resolver.Address{Addr: addr, ServerName: keyPrefix})
	} else {
		resp, _ := cli.Get(context.TODO(), keyPrefix, clientv3.WithPrefix())
		for v, _ := range resp.Kvs {
			addrList = append(addrList, resolver.Address{Addr: string(resp.Kvs[v].Value), ServerName: string(resp.Kvs[v].Key)})
		}
	}
	// 新版本etcd去除了NewAddress方法 以UpdateState代替
	r.cc.UpdateState(resolver.State{Addresses: addrList})
	WatchList(keyPrefix, r.ttl)
}

func (r etcdResolver) Scheme() string {
	return Schema
}

func (r etcdResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r etcdResolver) Close() {}
