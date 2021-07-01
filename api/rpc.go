package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"os"
	"wuyan94zl/services/etcd"
	"wuyan94zl/services/order/client"
)

// 获取服务
func getService(name string) (*grpc.ClientConn, error) {
	return etcd.GetService(client.ServerHost, name, client.TTL, etcd.GetRoundRobin())
}

func response(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": data,
		"msg":  "msg",
	})
}

func info(ctx *gin.Context) {
	conn, _ := getService(client.ServerName)
	defer conn.Close()
	cli := client.NewOrder(conn)
	res, err := cli.GetOrderInfo(context.Background(), &client.InfoRequest{Id: "1", UserId: "9999"})
	if err != nil {
		response(ctx, err.Error())
		return
	}
	data := make(map[string]interface{})
	json.Unmarshal([]byte(res.Data), &data)
	response(ctx, data)
}
func list(ctx *gin.Context) {
	conn, _ := getService(client.ServerName)
	defer conn.Close()
	cli := client.NewOrder(conn)
	res, err := cli.GetOrderList(context.Background(), &client.ListRequest{UserId: "8888"})
	if err != nil {
		response(ctx, err.Error())
		return
	}
	data := make(map[string]interface{})
	json.Unmarshal([]byte(res.Data), &data)
	response(ctx, data)
}

func main() {
	addr := ":8800"
	if len(os.Args) > 1 {
		addr = fmt.Sprintf(":%s", os.Args[1])
	}
	// api 服务
	ginHttp := gin.Default()
	ginHttp.GET("/info", info)
	ginHttp.GET("/list", list)
	ginHttp.Run(addr)

}
