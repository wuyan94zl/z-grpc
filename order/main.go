package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"wuyan94zl/services/etcd"
	"wuyan94zl/services/order/config"
	order "wuyan94zl/services/order/pb"
	"wuyan94zl/services/order/server"
)

func main() {

	if len(os.Args) > 1 {
		if port, e := strconv.Atoi(os.Args[1]); e == nil {
			config.Config.Port = port
		}
	}

	/** start etcd 注册、注销、健康监听 */
	// 服务地址
	ServiceAddr := fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port)

	// etcd 注册key
	etcdKey := config.Config.Rpc.Key

	// etcd Addr配置
	etcdAddr := config.Config.Rpc.Host

	// 注册服务
	go etcd.Register(etcdAddr, etcdKey, ServiceAddr, config.Config.TTL)

	// 监听服务 状态
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sig := <-ch
		fmt.Println("服务注销",etcdKey)
		// 服务注销
		etcd.UnRegister(etcdKey)
		if i, ok := sig.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()
	/** end etcd 注册、注销、健康监听 */

	/** start grpc 服务启动 */
	listen, err := net.Listen("tcp", ServiceAddr)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	// 实例化grpc Server
	s := grpc.NewServer()
	// 注册OrderServer
	order.RegisterOrderServer(s, server.OrderServer{})
	fmt.Println("Listen on " + ServiceAddr)
	s.Serve(listen)
	/** end grpc 服务启动 */
}
