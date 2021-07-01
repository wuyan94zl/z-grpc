package config

type config struct {
	Name string
	Host string
	Port int
	TTL  int64
	Rpc  rpc
}

type rpc struct {
	Host []string
	Key  string
}

var Config = &config{
	Name: "/order",
	Host: "0.0.0.0",
	Port: 8881,
	TTL:  5,
	Rpc: rpc{
		Host: []string{"127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379", "127.0.0.1:42379"},
		Key:  "/order.rpc",
	},
}
