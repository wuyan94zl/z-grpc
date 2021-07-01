package client

import (
	"context"
	"google.golang.org/grpc"
	"wuyan94zl/services/order/config"
	order "wuyan94zl/services/order/pb"
)

var ServerName = config.Config.Rpc.Key
var ServerHost = config.Config.Rpc.Host
var TTL = config.Config.TTL

type (
	InfoRequest   = order.InfoRequest
	ListRequest   = order.ListRequest
	OrderResponse = order.OrderResponse

	Order interface {
		GetOrderList(ctx context.Context, in *ListRequest) (*OrderResponse, error)
		GetOrderInfo(ctx context.Context, in *InfoRequest) (*OrderResponse, error)
	}

	defaultOrder struct {
		conn *grpc.ClientConn
	}
)

func NewOrder(cli *grpc.ClientConn) Order {
	return &defaultOrder{
		conn: cli,
	}
}

func (m *defaultOrder) GetOrderList(ctx context.Context, in *ListRequest) (*OrderResponse, error) {
	client := order.NewOrderClient(m.conn)
	return client.GetOrderList(ctx, in)
}
func (m *defaultOrder) GetOrderInfo(ctx context.Context, in *InfoRequest) (*OrderResponse, error) {
	client := order.NewOrderClient(m.conn)
	return client.GetOrderInfo(ctx, in)
}
