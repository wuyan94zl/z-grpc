##  学习 GRPC 微服务实践

### etcd 集群安装
```yml
version: '3'
networks:
  byfn:
  byfn1:

services:
  etcd1:
    image: quay.io/coreos/etcd:v3.5.0-rc.1
    container_name: etcd1
    command: etcd -name etcd1 -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379 -listen-peer-urls http://0.0.0.0:2380 -initial-cluster-token etcd-cluster -initial-cluster "etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380,etcd4=http://etcd4:2380" -initial-cluster-state new
    ports:
      - "12379:2379"
      - "12380:2380"
    networks:
      - byfn
  etcd2:
    image: quay.io/coreos/etcd:v3.5.0-rc.1
    container_name: etcd2
    command: etcd -name etcd2 -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379 -listen-peer-urls http://0.0.0.0:2380 -initial-cluster-token etcd-cluster -initial-cluster "etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380,etcd4=http://etcd4:2380" -initial-cluster-state new
    ports:
      - "22379:2379"
      - "22380:2380"
    networks:
      - byfn
  etcd3:
    image: quay.io/coreos/etcd:v3.5.0-rc.1
    container_name: etcd3
    command: etcd -name etcd3 -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379 -listen-peer-urls http://0.0.0.0:2380 -initial-cluster-token etcd-cluster -initial-cluster "etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380,etcd4=http://etcd4:2380" -initial-cluster-state new
    ports:
      - "32379:2379"
      - "32380:2380"
    networks:
      - byfn
  etcd4:
    image: quay.io/coreos/etcd:v3.5.0-rc.1
    container_name: etcd4
    command: etcd -name etcd4 -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379 -listen-peer-urls http://0.0.0.0:2380 -initial-cluster-token etcd-cluster -initial-cluster "etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380,etcd4=http://etcd4:2380" -initial-cluster-state new
    ports:
      - "42379:2379"
      - "42380:2380"
    networks:
      - byfn
```



### 运行多个服务
`go run order/main.go 8881` ：服务一  
`go run order/main.go 8882` ：服务二  
`go run order/main.go 8883` ：服务三  
服务提供2个接口，`info` 和 `list`  

### 运行api服务
`go run api/rpc.go`  
访问 `http://localhost:8807/info` 调用服务info接口    
访问 `http://localhost:8807/list` 调用服务list接口  
会以轮询的方式分别调用：服务一，服务二，服务三，返回结果中包含服务的地址端口信息  

结束上面三个服务之一，api调用中会自动剔除结束掉的服务，并不影响。（服务异常退出或死掉）  
增加服务四，api调用中会增加服务四的调用轮询。（达到随时扩容）

`protoc --go_out=plugins=grpc:. *.proto`
