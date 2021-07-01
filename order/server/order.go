package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"wuyan94zl/services/order/config"
	order "wuyan94zl/services/order/pb"
)

type OrderServer struct {

}

func (server OrderServer) GetOrderInfo(ctx context.Context, req *order.InfoRequest) (*order.OrderResponse, error) {
	data := make(map[string]string)
	data["id"] = req.Id
	data["user_id"] = req.UserId
	data["name"] = "订单名称"
	data["service"] = fmt.Sprintf("%s:%d",config.Config.Host,config.Config.Port)
	data["time"] = time.Now().String()
	str, err := json.Marshal(data)
	return &order.OrderResponse{Data: string(str)}, err
}

func (server OrderServer) GetOrderList(ctx context.Context, req *order.ListRequest) (*order.OrderResponse, error) {
	m := make(map[string]string)
	m["id"] = "1"
	m["user_id"] = req.UserId
	m["name"] = "订单名称1"
	var data []map[string]string
	data = append(data, m)
	m["id"] = "2"
	m["name"] = "订单名称2"
	data = append(data, m)
	rlt := make(map[string]interface{})
	rlt["status"] = 200
	rlt["service"] = fmt.Sprintf("%s:%d",config.Config.Host,config.Config.Port)
	rlt["data"] = data
	str, err := json.Marshal(rlt)
	return &order.OrderResponse{Data: string(str)}, err
}
