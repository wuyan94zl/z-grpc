package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

const Schema = "etcdService"

var cli *clientv3.Client

func Register(etcdAddr []string, name string, addr string, ttl int64) error {
	var err error
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   etcdAddr,
			DialTimeout: 15 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("err:%v", err)
		}
	}
	ticker := time.NewTicker(time.Second * time.Duration(ttl))
	go func() {
		for {
			_, err := cli.Get(context.Background(), "/"+Schema+name)
			if err != nil {
				log.Printf("Register:%s", err)
			} else {
				err = withAlive(name, addr, ttl)
				if err != nil {
					log.Printf("keep alive:%s", err)
				}
			}
			<-ticker.C
		}
	}()

	return nil
}

func withAlive(name string, addr string, ttl int64) error {
	leaseResp, err := cli.Grant(context.TODO(), ttl*2)
	if err != nil {
		return err
	}

	_, err = cli.Put(context.Background(), "/"+Schema+name, addr, clientv3.WithLease(leaseResp.ID))

	if err != nil {
		return fmt.Errorf("put etcd error:%v", err)
	}
	return nil
}

func UnRegister(name string) {
	if cli != nil {
		cli.Delete(context.Background(), "/"+Schema+name, clientv3.WithPrevKV())
	}
}

func GetService(config []string, name string, ttl int64, load Load) (*grpc.ClientConn, error) {
	newResolver := NewResolver(config, load, ttl)
	resolver.Register(newResolver)
	// 参数1 为服务地址，参数而为配置信息 当前标识负载均衡的轮询方式
	conn, err := grpc.Dial(newResolver.Scheme()+"://"+name, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("连接服务器失败:%v", err)
	}
	return conn, nil
}
